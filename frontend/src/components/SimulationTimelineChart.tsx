import { useMemo, useState } from "react";
import { useSimulator } from "../context/SimulatorContext";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

// Flatten the map into a single array for the chart
    interface ChartDataPoint {
      timestamp: number;
      cumulativeDamage: number;
      championName: string;
      championId: number;
      eventType?: string;
      eventData?: any; // Consider using a more specific type if the structure of eventData is known
}
    
// Custom tooltip component
const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;

    return (
      <div className="bg-background/95 border rounded-md p-2 shadow-md max-w-[400px] max-h-[300px] overflow-auto">
        <p className="font-medium">{`Time: ${label.toFixed(3)}s`}</p>
        <p>{`Champion: ${data.championName}`}</p>
        <p>{`Damage: ${data.cumulativeDamage.toFixed(2)}s`}</p>
        {data.eventType && (
          <>
            <p className="font-medium mt-2">Event Type:</p>
            <p className="text-xs">{data.eventType}</p>
          </>
        )}
        {data.eventData && (
          <>
            <p className="font-medium mt-2">Event Data:</p>
            <pre className="text-xs mt-1 whitespace-pre-wrap">
              {JSON.stringify(data.eventData, null, 2)}
            </pre>
          </>
        )}
      </div>
    );
  }
  return null;
};

// Custom dot component for events - Now aware of global hover state
const CustomDot = (props: any) => {
  // Destructure props, including payload and the new hoveredInfo/championColor
  const { cx, cy, payload, hoveredInfo, championColor } = props;

  // Only show dots for events with a type
  if (!payload.eventType) return null;

  // Determine if this specific dot is the one being hovered
  const isActive =
    hoveredInfo?.championId === payload.championId &&
    // Use a small tolerance for timestamp comparison due to potential float precision issues
    Math.abs(hoveredInfo?.timestamp - payload.timestamp) < 0.001;

  const eventType = payload.eventType;
  let baseFillColor = "#8884d8"; // Default color
  let baseRadius = 2;

  // Determine base appearance based on event type
  if (eventType.includes("DamageAppliedEvent")) {
    baseFillColor = "#ff6b6b"; // Red for damage
    baseRadius = 4;
  } else if (eventType.includes("AttackLandedEvent")) {
    baseFillColor = "#4ecdc4"; // Teal for attack landed
    baseRadius = 3;
  } else if (eventType.includes("SpellLandedEvent")) {
    baseFillColor = "#ffe66d"; // Yellow for spell landed
    baseRadius = 3;
  }
  // Add more conditions for other event types if needed

  // Determine final appearance based on isActive state
  const finalRadius = isActive ? 6 : baseRadius;
  // Use the line's color for fill when active, otherwise use the event-based color
  const finalFill = isActive ? championColor : baseFillColor;
  // Add a stroke only when active
  const finalStroke = isActive ? "#fff" : "none";
  const finalStrokeWidth = isActive ? 1 : 0;

  return (
    <circle
      cx={cx}
      cy={cy}
      r={finalRadius}
      fill={finalFill}
      stroke={finalStroke}
      strokeWidth={finalStrokeWidth}
      style={{ pointerEvents: "none" }} // Prevent dot from stealing hover events
    />
  );
};

const SimulationTimelineChart = () => {
  const { state } = useSimulator();
  const { simulationResults, simulationEvents, champions } = state;
  const [selectedChampionId, setSelectedChampionId] = useState<number | null>(
    null,
  );
  // --- Add Hover State ---
  const [hoveredInfo, setHoveredInfo] = useState<{
    championId: number;
    timestamp: number;
  } | null>(null);

  const chartData = useMemo(() => {
    if (!simulationResults || !simulationEvents || !champions) {
      return [];
    }

    // Create a map of entityId to champion name for easier lookup
    const entityToChampionMap = new Map();
    simulationResults.forEach((result) => {
      const champion = champions.find(
        (c) => c.apiName === result.championApiName,
      );
      entityToChampionMap.set(result.championEntityId, {
        name: champion?.name || result.championApiName,
        apiName: result.championApiName,
        id: result.championEntityId, // Add id here for easier lookup later
      });
    });

    // Process events to create chart data
    const championDamageData = new Map();

    // Initialize with starting points for each champion
    entityToChampionMap.forEach((champion, entityId) => {
      championDamageData.set(entityId, [
        {
          timestamp: 0,
          cumulativeDamage: 0,
          championName: champion.name,
          championId: entityId,
        },
      ]);
    });

    // Process all events
    simulationEvents.forEach((archivedEvent) => {
      const timestamp = archivedEvent.eventItem.Timestamp;
      const eventType = archivedEvent.eventType;

      // --- Determine the relevant entity ID based on event type ---
      let relevantEntityId: number | undefined = undefined;
      if (
        eventType.includes("DamageAppliedEvent") ||
        eventType.includes("AttackLandedEvent") ||
        eventType.includes("SpellLandedEvent")
      ) {
        relevantEntityId = archivedEvent.eventItem.Event.Source;
      } else if (archivedEvent.eventItem.Event.Entity !== undefined) {
        // Fallback for other potential event types that might use 'Entity'
        relevantEntityId = archivedEvent.eventItem.Event.Entity;
      }
      // --- End Entity ID determination ---

      // Skip if we couldn't find an entity or if it's not a tracked champion
      if (
        relevantEntityId === undefined ||
        !entityToChampionMap.has(relevantEntityId)
      ) {
        return;
      }

      const championName = entityToChampionMap.get(relevantEntityId).name;
      const championEvents = championDamageData.get(relevantEntityId) || [];
      const lastEvent = championEvents[championEvents.length - 1];

      // Handle damage events (updates cumulative damage)
      if (eventType.includes("DamageAppliedEvent")) {
        const damage = archivedEvent.eventItem.Event.FinalTotalDamage || 0;
        const newCumulativeDamage = (lastEvent?.cumulativeDamage || 0) + damage;

        championEvents.push({
          timestamp,
          cumulativeDamage: newCumulativeDamage,
          championName,
          championId: relevantEntityId,
          eventType,
          eventData: archivedEvent.eventItem.Event,
        });
      }
      // Handle other events (adds a dot without changing cumulative damage)
      else {
        if (championEvents.length > 0) {
          championEvents.push({
            ...lastEvent,
            timestamp,
            eventType,
            eventData: archivedEvent.eventItem.Event,
          });
        }
      }

      championDamageData.set(relevantEntityId, championEvents);
    });

    let chartDataPoints: ChartDataPoint[] = [];
    championDamageData.forEach((events) => {
      chartDataPoints = [...chartDataPoints, ...events];
    });

    // Sort by timestamp, and prioritize DamageAppliedEvent at same timestamp
    return chartDataPoints.sort((a, b) => {
      if (a.timestamp !== b.timestamp) {
        return a.timestamp - b.timestamp; // Primary sort: timestamp
      }

      // Secondary sort: Prioritize DamageAppliedEvent by putting it last
      const aIsDamage = a.eventType?.includes("DamageAppliedEvent");
      const bIsDamage = b.eventType?.includes("DamageAppliedEvent");

      if (aIsDamage && !bIsDamage) {
        return 1; // a (damage) comes after b (non-damage)
      }
      if (!aIsDamage && bIsDamage) {
        return -1; // b (damage) comes after a (non-damage)
      }

      // If both are damage or both are not, maintain relative order or sort by ID
      return a.championId - b.championId;
    });
    // --- End of existing chartData calculation ---
  }, [simulationResults, simulationEvents, champions]);

  // Group data by champion for separate lines
  const championGroups = useMemo(() => {
    if (chartData.length === 0) return [];

    const groups = new Map();
    // Use entityToChampionMap created in chartData calculation for consistency
    const entityToChampionMap = new Map();
    if (simulationResults && champions) {
      simulationResults.forEach((result) => {
        const champion = champions.find(
          (c) => c.apiName === result.championApiName,
        );
        entityToChampionMap.set(result.championEntityId, {
          id: result.championEntityId,
          name: champion?.name || result.championApiName,
          color: getChampionColor(result.championEntityId), // Assign color here
        });
      });
    }

    // Ensure all champions involved in the simulation are included, even if they have no events yet
    entityToChampionMap.forEach((championInfo) => {
      if (!groups.has(championInfo.id)) {
        groups.set(championInfo.id, championInfo);
      }
    });

    // Also iterate chartData to catch any potential edge cases (though entityToChampionMap should cover it)
    chartData.forEach((event) => {
      if (!groups.has(event.championId)) {
        groups.set(event.championId, {
          id: event.championId,
          name: event.championName,
          color: getChampionColor(event.championId),
        });
      }
    });

    return Array.from(groups.values());
  }, [chartData, simulationResults, champions]); // Add dependencies

  // Generate random color based on champion ID for consistency
  function getChampionColor(id: number) {
    // Simple hash function to generate color
    const hue = (id * 137) % 360;
    return `hsl(${hue}, 70%, 60%)`;
  }

  // Handler for clicking legend or line
  const handleChampionSelect = (championId: number) => {
    setSelectedChampionId((prevId) =>
      prevId === championId ? null : championId,
    );
  };

  if (!simulationResults || !simulationEvents || chartData.length === 0) {
    return (
      <div className="text-center text-muted-foreground italic py-8">
        Run a simulation to see culmulative damage over time.
      </div>
    );
  }

  // Filter champion groups based on selection
  const displayedChampionGroups =
    selectedChampionId === null
      ? championGroups
      : championGroups.filter((champion) => champion.id === selectedChampionId);

  // --- Tooltip Content Wrapper to Update State ---
  const renderTooltipContent = (tooltipProps: any) => {
    const { active, payload } = tooltipProps;

    if (active && payload && payload.length) {
      // Get data from the first payload item (since shared=false)
      const pointData = payload[0].payload;
      const currentHover = {
        championId: pointData.championId,
        timestamp: pointData.timestamp,
      };
      // Update state only if it changed to avoid potential re-renders
      if (
        hoveredInfo?.championId !== currentHover.championId ||
        hoveredInfo?.timestamp !== currentHover.timestamp
      ) {
        // Use rAF to batch state update slightly, might help with performance/flicker
        requestAnimationFrame(() => {
          setHoveredInfo(currentHover);
        });
      }
    } else {
      // Clear hover state if tooltip is inactive
      if (hoveredInfo !== null) {
        requestAnimationFrame(() => {
          setHoveredInfo(null);
        });
      }
    }

    // Render the actual tooltip content using the original component
    return <CustomTooltip {...tooltipProps} />;
  };

  return (
    <div className="w-full h-[400px] mt-4 mb-2">
      <p className="text-sm text-gray-500 mb-2">
        assumption: attack windup/wind-down takes 0s, cast windup/wind-down
        takes 1s each. used expected critically strike damage values for
        simulation.
      </p>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart
          data={chartData}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
          // Clear hover state if mouse leaves the entire chart area
          onMouseLeave={() => setHoveredInfo(null)}
        >
          <CartesianGrid strokeDasharray="3 3" opacity={0.3} />
          <XAxis
            dataKey="timestamp"
            label={{
              value: "Time (seconds)",
              position: "insideBottomRight",
              offset: -10,
            }}
            domain={[0, 30]}
            ticks={[0, 5, 10, 15, 20, 25, 30]}
            type="number"
          />
          <YAxis
            label={{
              value: "Cumulative Damage",
              angle: -90,
              position: "insideLeft",
            }}
            // Dynamically adjust domain if a champion is selected
            domain={selectedChampionId !== null ? ["auto", "auto"] : undefined}
            allowDataOverflow={true} // Prevent clipping when domain changes
          />
          {/* --- Use Tooltip Wrapper --- */}
          <Tooltip
            content={renderTooltipContent}
            shared={false}
            // Optional: Adjust position if needed
            // position={{ y: 0 }}
          />
          <Legend
            onClick={(data) => {
              let championIdToSelect: number | null = null;

              // Assert the type of data.payload to include the 'data' property.
              // We expect data.payload to be the props passed to the <Line /> component,
              // which includes a 'data' prop of type ChartDataPoint[].
              const typedPayload = data.payload as { data?: ChartDataPoint[] };

              if (
                data && // data is the legend item descriptor from Recharts
                typedPayload && // Check the (now asserted) payload
                Array.isArray(typedPayload.data) && // Access 'data' via the typedPayload
                typedPayload.data.length > 0 &&
                typedPayload.data[0] &&
                typeof typedPayload.data[0].championId === "number"
              ) {
                championIdToSelect = typedPayload.data[0].championId;
              } else {
                // Fallback logic: if the expected path is not valid, try finding by color
                console.warn(
                  "Legend click: Could not find championId in data.payload.data. Payload:",
                  data.payload, // Log the original payload for debugging if the primary path fails
                  "Falling back to color-based lookup.",
                );
                const clickedChampion = championGroups.find(
                  (c) => c.color === data.color, // data.color is available on the legend item descriptor
                );
                if (clickedChampion) {
                  championIdToSelect = clickedChampion.id;
                } else {
                  console.error(
                    "Legend click: Could not find champion by data path or color.",
                    data, // Log the entire legend item descriptor for debugging
                  );
                }
              }

              if (championIdToSelect !== null) {
                handleChampionSelect(championIdToSelect);
              }
            }}
            wrapperStyle={{ cursor: "pointer" }} // Add pointer cursor to legend
          />
          {displayedChampionGroups.map((champion) => (
            <Line
              key={champion.id}
              type="linear"
              dataKey="cumulativeDamage"
              data={chartData.filter((d) => d.championId === champion.id)}
              name={champion.name}
              stroke={champion.color}
              strokeWidth={2}
              // --- Pass hoverInfo and championColor to CustomDot via dot prop ---
              dot={(dotProps: any) => {
                const { key, ...restDotProps } = dotProps; // Destructure key and rest of the props
                return (
                  <CustomDot
                    key={key} // Pass key directly
                    {...restDotProps} // Spread the remaining props
                    hoveredInfo={hoveredInfo}
                    championColor={champion.color} // Pass the specific line's color
                  />
                );
              }}
              // --- Remove activeDot prop ---
              activeDot={false} // Disable default activeDot behavior
              connectNulls
              isAnimationActive={false}
              onClick={() => handleChampionSelect(champion.id)}
              style={{ cursor: "pointer" }}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default SimulationTimelineChart;
