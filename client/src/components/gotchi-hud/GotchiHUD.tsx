import { useEffect, useMemo, useState } from "react";
import "./GotchiHUD.css";
import { fetchSingleGotchiSVGs, GotchiSVGSet } from "../../utils/gotchi-loader/FetchGotchis";
import { emotiveAccessories } from "./emotive-accessories";
import { Sprite } from "../Sprite";
import { SpeechBubble } from "../SpeechBubble";
import StatBar from "../menu/StatBar";

interface Props {}

const animationSequence = ["normal", "happy", "normal", "sad", "normal", "mad"] as const;
type AnimationKey = typeof animationSequence[number];

export const GotchiHUD = (props: Props) => {
  const [gotchiSvgSet, setGotchiSvgSet] = useState<GotchiSVGSet | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [animationIndex, setAnimationIndex] = useState(0);

    const [ecto, setEcto] = useState(500);
    const [spark, setSpark] = useState(500);
    const [pulse, setPulse] = useState(500);

  // Fetch gotchi on mount
  useEffect(() => {
    const fetchGotchi = async () => {
      try {
        const gotchiSvgSet = await fetchSingleGotchiSVGs("115", 256);
        setGotchiSvgSet(gotchiSvgSet);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch gotchi data for 115: " + err);
        setLoading(false);
      }
    };

    fetchGotchi();
  }, []);

  // Animation cycle
  useEffect(() => {
    const interval = setInterval(() => {
      setAnimationIndex((prev) => (prev + 1) % animationSequence.length);
    }, 500);

    return () => clearInterval(interval);
  }, []);

  if (loading || !gotchiSvgSet?.anims) {
    return <div className="gotchi-hud">Loading...</div>;
  }

  if (error) {
    return <div className="gotchi-hud">{error}</div>;
  }

  const currentAnim = animationSequence[animationIndex];
  const svgHtml = gotchiSvgSet.anims[currentAnim] ?? gotchiSvgSet?.front ?? "";
  const emotionFrame = emotiveAccessories[currentAnim];

  return (
    <div className="gotchi-hud">
      <div className="gotchi-speech">
        <SpeechBubble
          text="Hello, World! What is happening today???"
          width={256-48}
        />
      </div>
      <div
        className="gotchi-svg"
        style={{ animationDuration: "0.5s" }}
        dangerouslySetInnerHTML={{ __html: svgHtml }}
      />
      <div className="gotchi-emotion">
        <Sprite
          frameIndex={emotionFrame}
          frameWidth={48}
          frameHeight={48}
          sheetColumns={8}
          imageUrl="/assets/emoticons/emoticons_48px.png"
          frameMargin={2}
          frameSpacing={4}
        />
      </div>
      <div className="gotchi-action">
        <Sprite
          frameIndex={emotionFrame}
          frameWidth={48}
          frameHeight={48}
          sheetColumns={8}
          imageUrl="/assets/emoticons/emoticons_48px.png"
          frameMargin={2}
          frameSpacing={4}
        />
      </div>
      <div className="gotchi-esp">
        <StatBar label="Ecto" value={ecto} max={1000} color="#CA52C9" />
        <StatBar label="Spark" value={spark} max={1000} color="#0098DC" />
        <StatBar label="Pulse" value={pulse} max={1000} color="#F5555D" />
      </div>
    </div>
  );
};
