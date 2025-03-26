// HoverInfo.tsx
import { useEffect, useState } from "react";
import "./HoverInfo.css";

export const HoverInfo = () => {
    const [hoverData, setHoverData] = useState<any | null>(null);

    useEffect(() => {
        const handleHover = (event: Event) => {
            const customEvent = event as CustomEvent;
            // console.log("hover: ", customEvent.detail);
            setHoverData(customEvent.detail);
        };

        window.addEventListener("entityHover", handleHover);

        return () => {
            window.removeEventListener("entityHover", handleHover);
        };
    }, []);

    if (!hoverData) return null;

    return (
        <div className="hover-info-container">
            <h3 className="hover-info-title">Entity Info</h3>
            <pre className="hover-info-content">{JSON.stringify(hoverData, null, 2)}</pre>
        </div>
    );

};
