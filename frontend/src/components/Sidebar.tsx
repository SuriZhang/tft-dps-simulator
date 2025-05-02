import React, { useState } from 'react';
import { cn } from '../lib/utils';
import { Button } from './ui/button'; // Import Button
import { ChevronLeft, ChevronRight } from 'lucide-react'; // Import icons

interface SidebarProps {
  children: React.ReactNode;
  title: string;
  position: 'left' | 'right';
  defaultCollapsed?: boolean;
}

const Sidebar: React.FC<SidebarProps> = ({
  children,
  title,
  position,
  defaultCollapsed = false
}) => {
  const [collapsed, setCollapsed] = useState(defaultCollapsed);

  return (
    <div
      className={cn(
        "bg-card rounded-lg shadow-lg transition-all duration-300 sidebar-transition",
        collapsed ? "w-12" : "w-64",
        position === 'left' ? "mr-4" : "ml-4"
      )}
    >
      {/* Header */}
      <div
        className={cn(
          "flex items-center p-4 border-b border-gray-800 cursor-pointer",
          collapsed ? "justify-center" : "justify-between"
        )}
        onClick={() => setCollapsed(!collapsed)}
      >
        {!collapsed && (
          <h2 className="text-lg font-bold text-white">{title}</h2>
        )}
        {/* Use Button component with lucide icons */}
        <Button variant="ghost" size="icon" className="w-6 h-6 text-gray-400 hover:text-white">
          {collapsed ? (
            position === 'left' ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />
          ) : (
            position === 'left' ? <ChevronLeft className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />
          )}
          <span className="sr-only">{collapsed ? 'Expand' : 'Collapse'} Sidebar</span>
        </Button>
      </div>

      {/* Content */}
      <div className={cn("overflow-hidden", collapsed && "opacity-0")}>
        {!collapsed && (
          <div className="p-4">
            {children}
          </div>
        )}
      </div>
    </div>
  );
};

export default Sidebar;
