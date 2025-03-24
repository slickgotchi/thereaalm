interface Props {
    scene: Phaser.Scene;
    x: number;
    y: number;
    currentHP: number;
    maxHP: number;
    trackingSprite: Phaser.GameObjects.Sprite;
}


export class HPBar {
    private scene: Phaser.Scene;
    private maxHP: number;
    private currentHP: number;
    private width: number = 42;  // Total width of the background
    private height: number = 6;  // Total height of the background
    private fillWidth: number = 40;  // Width of the fill bar when full
    private fillHeight: number = 4;  // Height of the fill bar
    private background: Phaser.GameObjects.Rectangle;
    private fill: Phaser.GameObjects.Rectangle;
    private trackingSprite: Phaser.GameObjects.Sprite;

    constructor(props: Props) {
        const {scene, x, y, currentHP, maxHP, trackingSprite} = props;

        this.trackingSprite = trackingSprite;
        this.scene = scene;
        this.maxHP = maxHP;
        this.currentHP = currentHP;

        // Create the black background rectangle (42x8)
        this.background = this.scene.add.rectangle(
            x,
            y,
            this.width,  // 42 pixels wide
            this.height, // 8 pixels high
            0x000000     // Black color
        );
        this.background.setOrigin(0, 0);

        // Create the green fill bar (40x6 when full, scaled by HP ratio)
        this.fill = this.scene.add.rectangle(
            x,
            y,
            this.fillWidth * (this.currentHP / this.maxHP), // Scale width based on HP ratio
            this.fillHeight,                         // 6 pixels high
            0x00ff00                                 // Green color
        );
        this.fill.setOrigin(0, 0);

        // Ensure the fill bar is on top of the background
        this.background.setDepth(5000);
        this.fill.setDepth(5001);

        this.setPosition(x,y);
    }

    // Update the HP bar when HP changes
    public updateHP(currentHP: number): void {
        if (!this.trackingSprite) {
            this.destroy();
            return;
        }

        this.currentHP = Math.max(0, Math.min(currentHP, this.maxHP)); // Clamp between 0 and maxHP
        this.fill.setScale(this.currentHP / this.maxHP, 1); // Adjust width of the fill bar

        this.setPosition(this.trackingSprite.x, this.trackingSprite.y);
    }

    // Move the HP bar to a new position (e.g., follow the character)
    public setPosition(x: number, y: number): void {
        const offsetX = (64-this.width) / 2;

        this.background.setPosition(x + offsetX, y+1);
        this.fill.setPosition(x + offsetX + 1, y+1+1);
    }

    // Destroy the HP bar when no longer needed
    public destroy(): void {
        this.background.destroy();
        this.fill.destroy();
    }
}