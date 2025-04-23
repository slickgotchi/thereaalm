import React from "react";
import "./SpeechBubble.css";

interface SpeechBubbleProps {
  text: string;
  width: number;
}

export const SpeechBubble: React.FC<SpeechBubbleProps> = ({ text, width }) => {
  return (
    <div className="speech-bubble" style={{ width: `${width}px` }}>
      <div className="speech-bubble-inner">
        {text}
      </div>
    </div>
  );
};
