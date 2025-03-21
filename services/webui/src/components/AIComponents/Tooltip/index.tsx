import React, { useState } from "react";

interface ToolTipProps {
    text: string;
    children: React.ReactNode;
    className?: string;
    position?: string;
}

const Tooltip = ({ text, children,className,position }:ToolTipProps) => {
  const [isVisible, setIsVisible] = useState(false);
  const checkPosition =()=>{
    switch(position){
      case 'top-left':
        return 'left-0 top-0'
      case 'top-right':
        return 'right-0 top-0'
      case 'bottom-right':
        return 'right-0'
      case 'bottom-left':
        return '-left-full';
      default:
        return 'left-0'
    }
  }
  return (
    <div
      className={`relative sm:inline-block  ${className}`}
      onMouseEnter={() => setIsVisible(true)}
      onMouseLeave={() => setIsVisible(false)}
    >
      {isVisible && (
        <div className={`absolute bottom-full sm:inline-block hidden   ${position ? checkPosition() : 'left-0'} transform mb-2  w-max max-w-[20vw] px-3 py-2 text-sm text-white bg-slate-600 rounded-2xl shadow-lg`}>
          {text}
          {/* <div className="absolute top-full left-12 transform  border-8 border-transparent border-t-gray-900"></div> */}
        </div>
      )}
      {children}
    </div>
  );
};

export default Tooltip;
