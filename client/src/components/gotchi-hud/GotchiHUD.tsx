import { useEffect, useMemo, useState } from "react";
import "./GotchiHUD.css";
import { fetchSingleGotchiSVGs, GotchiSVGSet } from "../../utils/gotchi-loader/FetchGotchis";
import { actionFrameNumbers, emotionFrameNumbers } from "./icon-accessories";
import { Sprite } from "../Sprite";
import { SpeechBubble } from "../SpeechBubble";
import StatBar from "../menu/StatBar";
import { GotchiEntity } from "../../phaser/entities/GotchiEntity";

interface Props {
  selectedGotchiEntity: { type: string; [key: string]: any } | null;
}

const animationSequence = ["normal", "happy", "normal", "sad", "normal", "mad"] as const;
type AnimationKey = typeof animationSequence[number];

export const GotchiHUD: React.FC<Props> = ({selectedGotchiEntity}) => {
  const [gotchiSvgSet, setGotchiSvgSet] = useState<GotchiSVGSet | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [animationIndex, setAnimationIndex] = useState(0);

  const [gotchiId, setGotchiId] = useState<string | null>(null);

  // ESP stats
  const [ecto, setEcto] = useState(500);
  const [spark, setSpark] = useState(500);
  const [pulse, setPulse] = useState(500);

  // action
  const [action, setAction] = useState("roam");

  // emotion
  const [emotion, setEmotion] = useState("happy");

  // check for gotchi entity selection
  useEffect(() => {
    if (!selectedGotchiEntity || selectedGotchiEntity.type !== "gotchi") return;
    
    // function to be called when polling
    const fetchGotchiData = async () => {
      
      const newGotchiId = selectedGotchiEntity?.data?.gotchiId;
      if (!newGotchiId) return;
      
      setGotchiId(newGotchiId);

      if (!gotchiId) return;

      const gotchiEntity = GotchiEntity.activeGotchis.get(gotchiId);

      // set ESP stats
      setEcto(gotchiEntity?.data?.stats?.ecto || 0);
      setSpark(gotchiEntity?.data?.stats?.spark || 0);
      setPulse(gotchiEntity?.data?.stats?.pulse || 0);

      // set action
      const currentAction = gotchiEntity?.data?.actionPlan.currentAction;
      if (currentAction) {
        console.log("actionPlan: ", gotchiEntity?.data?.actionPlan);
        setAction(currentAction.type);
      }

      // set emotion
      setEmotion("happy");
    };

    fetchGotchiData();

    const intervalId = setInterval(fetchGotchiData, 3000);

    return () => clearInterval(intervalId);

  }, [selectedGotchiEntity, gotchiId])

  useEffect(() => {
    
    const loadSVG = async () => {
      if (!gotchiId) return;

      // load svg
      try {
        const gotchiSvgSet = await fetchSingleGotchiSVGs(gotchiId, 256);
        setGotchiSvgSet(gotchiSvgSet);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch gotchi data for 115: " + err);
        setLoading(false);
      }
    }

    loadSVG();

  }, [gotchiId])

  // Animation cycle
  useEffect(() => {
    const interval = setInterval(() => {
      setAnimationIndex((prev) => prev === 1 ? 0 : 1);
    }, 500);

    return () => clearInterval(interval);
  }, []);

  if (loading || !gotchiSvgSet?.anims) {
    return <div className="gotchi-hud">Loading...</div>;
  }

  if (error) {
    return <div className="gotchi-hud">{error}</div>;
  }

  // Safely access anims value by string name
  const getAnimSvg = (emotionKey: string): string | undefined => {
    // Ensure anims exists and the key is valid
    if (gotchiSvgSet.anims && emotionKey in gotchiSvgSet.anims) {
      return gotchiSvgSet.anims[emotionKey as keyof typeof gotchiSvgSet.anims];
    }
    return undefined; // Fallback if anims or key is missing
  };

  const svgHtml = animationIndex === 0 ? getAnimSvg("normal") : getAnimSvg(emotion);
  if (!svgHtml) return;
  
  const emotionFrame = emotionFrameNumbers[emotion];
  const actionFrame = actionFrameNumbers[action];

  return (
    <div className="gotchi-hud">
      <div className="gotchi-speech">
        <SpeechBubble
          text="Hello, World! What's happening today?"
          width={256+64}
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
          frameIndex={actionFrame}
          frameWidth={48}
          frameHeight={48}
          sheetColumns={8}
          imageUrl="/assets/spritesheets/actionicons_spritesheet.png"
          frameMargin={0}
          frameSpacing={0}
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
