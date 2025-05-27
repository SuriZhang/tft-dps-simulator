import React from "react";
import { useSimulator } from "../context/SimulatorContext";
import { Button } from "./ui/button"; 
import {
  Trash2,
  Settings,
  Share2,
  Copy,
  UploadCloud,
} from "lucide-react"; 
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "./ui/tooltip";

const ControlBar: React.FC = () => {
  const { dispatch } = useSimulator();

  // Toggle states for UI controls - Assuming these might come from context later
  const [showNames, setShowNames] = React.useState(true);
  const [useSkins, setUseSkins] = React.useState(true);
  const [mouseHoverInfo, setMouseHoverInfo] = React.useState(true);
  const [positioningMode, setPositioningMode] = React.useState(false);

  // Handle clear board action
  const handleClearBoard = () => {
    dispatch({ type: "CLEAR_BOARD" });
  };

  // Placeholder handlers
  const handleShare = () => console.log("Share TBD");
  const handleCopyCode = () => console.log("Copy Code TBD");
  const handleImportCode = () => console.log("Import Code TBD");
  const handleSettings = () => console.log("Settings TBD");

  // Update global state when switch changes (if applicable)
  const handleUseSkinsChange = (checked: boolean) => {
    setUseSkins(checked);
    // dispatch({ type: 'SET_USE_SKINS', payload: checked });
  };

  const handleMouseHoverInfoChange = (checked: boolean) => {
    setMouseHoverInfo(checked);
    // dispatch({ type: 'SET_MOUSE_HOVER_INFO', payload: checked });
  };

  const handleShowNamesChange = (checked: boolean) => {
    setShowNames(checked);
    // dispatch({ type: 'SET_SHOW_NAMES', payload: checked });
  };

  const handlePositioningModeChange = (checked: boolean) => {
    setPositioningMode(checked);
    // dispatch({ type: 'SET_POSITIONING_MODE', payload: checked });
  };

  return (
    <div className="bg-card rounded-lg shadow-lg p-4 mb-4">
      <div className="flex flex-wrap justify-between items-center gap-4">
        {/* Set info */}
        <div className="flex items-center space-x-4">
          {/* Consider using Badge component here if appropriate */}
          <div className="bg-accent rounded-md px-3 py-1">
            <span className="text-white font-medium">SET 14</span>
          </div>
          <div className="text-gray-400 text-sm">
            SET 13 {/* Potentially a Button or Link */}
          </div>
          <div className="text-sm text-gray-300">
            Right click a unit on board to mark it as 3-star.
          </div>

          <TooltipProvider delayDuration={100}>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="outline" size="sm" onClick={handleCopyCode}>
                  <Copy className="h-4 w-4 mr-1" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                Copy Code
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <TooltipProvider delayDuration={100}>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="outline" size="sm" onClick={handleImportCode}>
            <UploadCloud className="h-4 w-4 mr-1" /> 
          </Button>
              </TooltipTrigger>
              <TooltipContent>
                Import Code
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <TooltipProvider delayDuration={100}>
            <Tooltip>
              <TooltipTrigger asChild>
                 <Button variant="outline" size="sm" onClick={handleShare}>
            <Share2 className="h-4 w-4 mr-1" />
          </Button>
              </TooltipTrigger>
              <TooltipContent>
                Share Code
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <Button variant="destructive" size="sm" onClick={handleClearBoard}>
            <Trash2 className="h-4 w-4 mr-1" /> CLEAR BOARD
          </Button>
        </div>

        {/* Toggle controls */}
        {/* <div className="flex items-center space-x-4 ml-auto">
          <div className="flex items-center space-x-2">
            <Switch
              id="mouse-hover-switch"
              checked={mouseHoverInfo}
              onCheckedChange={handleMouseHoverInfoChange}
              aria-label="Enable mouse hover information"
            />
            <Label
              htmlFor="mouse-hover-switch"
              className="text-sm text-gray-300 cursor-pointer"
            >
              Hover Info
            </Label>
          </div>
          <div className="flex items-center space-x-2">
            <Switch
              id="show-names-switch"
              checked={showNames}
              onCheckedChange={handleShowNamesChange}
              aria-label="Show champion names"
            />
            <Label
              htmlFor="show-names-switch"
              className="text-sm text-gray-300 cursor-pointer"
            >
              Show Names
            </Label>
          </div>
          <div className="flex items-center space-x-2">
            <Switch
              id="positioning-mode-switch"
              checked={positioningMode}
              onCheckedChange={handlePositioningModeChange}
              aria-label="Enable positioning mode"
            />
            <Label
              htmlFor="positioning-mode-switch"
              className="text-sm text-gray-300 cursor-pointer"
            >
              Positioning Mode
            </Label>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={handleSettings}
            title="Settings"
          >
            <Settings className="h-5 w-5 text-gray-400" />
          </Button>
        </div> */}
      </div>
    </div>
  );
};

export default ControlBar;
