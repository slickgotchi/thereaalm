import { useEffect, useState } from "react";
import "./HoverInfo.css";
import { GotchiHoverInfo } from "./GotchiHoverInfo";

interface HoverData {
  id: string;
  type: string;
  data: any;
}

export const HoverInfo = () => {
  const [selectedEntity, setSelectedEntity] = useState<HoverData | null>(null);

  useEffect(() => {
    const handleSelection = (event: Event) => {
        const customEvent = event as CustomEvent;
        const newData = customEvent.detail;
        console.log('received event: ', newData);
      setSelectedEntity(newData); // Set to null or new data directly
    };

    window.addEventListener("entitySelection", handleSelection);

    return () => {
      window.removeEventListener("entitySelection", handleSelection);
    };
  }, []);

  if (!selectedEntity) return null;

  return (
    <div className="hover-info-container">
        <div className="hover-info-fallback">
          <h3 className="hover-info-title">Entity Info</h3>
          <pre className="hover-info-content">
            {JSON.stringify(selectedEntity, null, 2)}
          </pre>
        </div>
    </div>
  );
};

/*
import { useEffect, useState } from "react";
import "./HoverInfo.css";
import { GotchiHoverInfo } from "./GotchiHoverInfo";

interface HoverData {
  type: string;
  data: any;
}

export const HoverInfo = () => {
  const [hoverData, setHoverData] = useState<HoverData | null>(null);

  useEffect(() => {
    const handleHover = (event: Event) => {
      const customEvent = event as CustomEvent;
      const newData = customEvent.detail;

      if (!newData || !newData.data) {
        setHoverData(null); // No data, set to null and show nothing
      } else {
        setHoverData(newData); // Set new hover data when available
      }
    };

    window.addEventListener("entityHover", handleHover);

    // Cleanup on component unmount
    return () => {
      window.removeEventListener("entityHover", handleHover);
    };
  }, []);

  // Only render if there's hoverData (no need for a loading screen)
  if (!hoverData) return null;

  return (
    <div className="hover-info-container">
        <div className="hover-info-fallback">
        <h3 className="hover-info-title">Entity Info</h3>
        <pre className="hover-info-content">
            {JSON.stringify(hoverData, null, 2)}
        </pre>
        </div>
    </div>
  );

//   return (
//     <div className="hover-info-container">
//       {(() => {
//         switch (hoverData?.type) {
//           case "gotchi":
//             return <GotchiHoverInfo data={hoverData?.data} />;
//           default:
//             return (
//               <div className="hover-info-fallback">
//                 <h3 className="hover-info-title">Entity Info</h3>
//                 <pre className="hover-info-content">
//                   {JSON.stringify(hoverData, null, 2)}
//                 </pre>
//               </div>
//             );
//         }
//       })()}
//     </div>
//   );
};
*/