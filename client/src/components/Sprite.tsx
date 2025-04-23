interface SpriteProps {
    frameIndex: number;
    frameWidth: number;
    frameHeight: number;
    imageUrl: string;
    sheetColumns: number;     // how many columns in the spritesheet
    frameMargin?: number;     // margin around the entire sheet
    frameSpacing?: number;    // spacing between frames
  }
  
  export const Sprite = ({
    frameIndex,
    frameWidth,
    frameHeight,
    imageUrl,
    sheetColumns,
    frameMargin = 0,
    frameSpacing = 0,
  }: SpriteProps) => {
    const column = frameIndex % sheetColumns;
    const row = Math.floor(frameIndex / sheetColumns);
  
    const offsetX = -(frameMargin + column * (frameWidth + frameSpacing));
    const offsetY = -(frameMargin + row * (frameHeight + frameSpacing));
  
    return (
      <div
        style={{
          width: `${frameWidth}px`,
          height: `${frameHeight}px`,
          backgroundImage: `url(${imageUrl})`,
          backgroundPosition: `${offsetX}px ${offsetY}px`,
          backgroundRepeat: 'no-repeat',
          imageRendering: 'pixelated',
        }}
      />
    );
  };
  