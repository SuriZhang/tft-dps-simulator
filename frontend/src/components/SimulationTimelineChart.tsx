import React, { useMemo } from "react";
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

// Custom tooltip component
const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    
    return (
      <div className="bg-background/95 border rounded-md p-2 shadow-md max-w-[400px] max-h-[300px] overflow-auto">
        <p className="font-medium">{`Time: ${label.toFixed(2)}s`}</p>
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

// Custom dot component for events
const CustomDot = (props: any) => {
  const { cx, cy, payload } = props;
  
  // Only show dots for events with a type
  if (!payload.eventType) return null;

  const eventType = payload.eventType;
  let fillColor = "#8884d8"; // Default color (e.g., purple)
  let radius = 2;

  if (eventType.includes("DamageAppliedEvent")) {
    fillColor = "#ff6b6b"; // Red for damage
    radius = 4;
  }
  if (eventType.includes("AttackLandedEvent")) {
    fillColor = "#4ecdc4"; // Teal for attack landed
    radius = 4;
  }
  if (eventType.includes("SpellLandedEvent")) {
    fillColor = "#ffe66d"; // Yellow for spell landed
    radius = 4;
  }
  
  return (
    <circle 
      cx={cx} 
      cy={cy} 
      r={radius}
      fill={fillColor}
      stroke="none"
    />
  );
};

const SimulationTimelineChart = () => {
  const { state } = useSimulator();
  const { simulationResults, simulationEvents, champions } = state;

  const chartData = useMemo(() => {
    if (!simulationResults || !simulationEvents || !champions) {
      return [];
    }

    // Create a map of entityId to champion name for easier lookup
    const entityToChampionMap = new Map();
    simulationResults.forEach((result) => {
      const champion = champions.find(c => c.apiName === result.championApiName);
      entityToChampionMap.set(result.championEntityId, {
        name: champion?.name || result.championApiName,
        apiName: result.championApiName,
      });
    });

    // Process events to create chart data
    const championDamageData = new Map();

    // Initialize with starting points for each champion
    entityToChampionMap.forEach((champion, entityId) => {
      championDamageData.set(entityId, [{
        timestamp: 0,
        cumulativeDamage: 0,
        championName: champion.name,
        championId: entityId,
      }]);
    });

    // Process all events
    simulationEvents.forEach((archivedEvent) => {
      const timestamp = archivedEvent.eventItem.Timestamp;
      const eventType = archivedEvent.eventType;

      // --- Determine the relevant entity ID based on event type ---
      let relevantEntityId: number | undefined = undefined;
      if (eventType.includes("DamageAppliedEvent") || 
          eventType.includes("AttackLandedEvent") || 
          eventType.includes("SpellLandedEvent")) {
        relevantEntityId = archivedEvent.eventItem.Event.Source;
      } else if (archivedEvent.eventItem.Event.Entity !== undefined) {
        // Fallback for other potential event types that might use 'Entity'
        relevantEntityId = archivedEvent.eventItem.Event.Entity;
      }
      // --- End Entity ID determination ---

      // Skip if we couldn't find an entity or if it's not a tracked champion
      if (relevantEntityId === undefined || !entityToChampionMap.has(relevantEntityId)) {
        return; 
      }

      // Skip specific event types we don't want to plot dots for
      if (eventType.includes("AttackCooldown") || eventType.includes("ChampionActionEvent")) {
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
          eventData: archivedEvent.eventItem.Event
        });
      }
      // Handle other events (adds a dot without changing cumulative damage)
      else {
         if (championEvents.length > 0) { // Ensure there's a previous event to copy damage from
           championEvents.push({
             ...lastEvent, // Copy previous data (especially cumulativeDamage)
             timestamp,    // Update timestamp
             eventType,    // Set the correct event type
             eventData: archivedEvent.eventItem.Event // Set the event data
           });
         }
         // Optional: Handle case where the very first event for a champion is not a damage event
         // else { /* Add logic if needed, e.g., push with cumulativeDamage: 0 */ }
      }
      
      championDamageData.set(relevantEntityId, championEvents);
    });

    // Flatten the map into a single array for the chart
    interface ChartDataPoint {
      timestamp: number;
      cumulativeDamage: number;
      championName: string;
      championId: number;
      eventType?: string;
      eventData?: any; // Consider using a more specific type if the structure of eventData is known
      }
      
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
  }, [simulationResults, simulationEvents, champions]);

  // Group data by champion for separate lines
  const championGroups = useMemo(() => {
    if (chartData.length === 0) return [];
    
    const groups = new Map();
    chartData.forEach(event => {
      if (!groups.has(event.championId)) {
        groups.set(event.championId, {
          id: event.championId,
          name: event.championName,
          color: getChampionColor(event.championId)
        });
      }
    });
    
    return Array.from(groups.values());
  }, [chartData]);

  // Generate random color based on champion ID for consistency
  function getChampionColor(id: number) {
    // Simple hash function to generate color
    const hue = (id * 137) % 360;
    return `hsl(${hue}, 70%, 60%)`;
  }

  if (!simulationResults || !simulationEvents || chartData.length === 0) {
    return <div className="p-4 text-center">Run a simulation to see damage timeline</div>;
  }

  return (
    <div className="w-full h-[400px] mt-4">
      <h2 className="text-lg font-medium mb-4">Damage Timeline</h2>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart
          data={chartData}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray="3 3" opacity={0.3} />
          <XAxis 
            dataKey="timestamp" 
            label={{ value: 'Time (seconds)', position: 'insideBottomRight', offset: -10 }}
            domain={[0, 30]}
            ticks={[0, 5, 10, 15, 20, 25, 30]}
            type="number"
          />
          <YAxis 
            label={{ value: 'Cumulative Damage', angle: -90, position: 'insideLeft' }}
          />
          <Tooltip content={<CustomTooltip />} />
          <Legend />
          
          {championGroups.map(champion => (
            <Line
              key={champion.id}
              type="linear"
              dataKey="cumulativeDamage"
              data={chartData.filter(d => d.championId === champion.id)}
              name={champion.name}
              stroke={champion.color}
              strokeWidth={2}
              dot={<CustomDot />}
              activeDot={{ r: 6 }}
              connectNulls
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default SimulationTimelineChart;