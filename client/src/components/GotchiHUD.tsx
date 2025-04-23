import { useEffect, useMemo, useState } from "react";
import "./GotchiHUD.css";
import { Aavegotchi, fetchSingleGotchiSVGs, GotchiSVGSet } from "../utils/FetchGotchis";


interface Props {
  
}

export const GotchiHUD = (props: Props) => {

  const [gotchiSvgSet, setGotchiSvgSet] = useState<GotchiSVGSet | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // fetch gotchi on mount
  useEffect(() => {
    const fetchGotchi = async () => {
      try {
        const gotchiSvgSet = await fetchSingleGotchiSVGs("115", 256);
        setGotchiSvgSet(gotchiSvgSet);
        setLoading(false);

      } catch (err) {
        setError("Failed to fetch gotchi data for 115" + err);
        setLoading(false);
      }
    }

    fetchGotchi();
  }, [])

  // Render loading, error, or SVG
  if (loading) {
    return <div className="gotchi-hud">Loading...</div>;
  }

  if (error) {
    return <div className="gotchi-hud">{error}</div>;
  }

  return (
    <div className="gotchi-hud">
      <div
        className="gotchi-svg"
        dangerouslySetInnerHTML={{ __html: gotchiSvgSet?.front || "" }}
      />
    </div>
  );
};
