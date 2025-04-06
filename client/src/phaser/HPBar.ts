interface Props {
    scene: Phaser.Scene;
    x: number;
    y: number;
    currentHP: number;
    maxHP: number;
    trackingSprite: Phaser.GameObjects.Sprite;
}

// HP bars are scaled on an assumed max 1000 pulse
export class HPBar {
    private scene: Phaser.Scene;
    private maxHP: number;
    private currentHP: number;
    private maxFillWidth: number = 32;
    private fillHeight: number = 4;
    private padding: number = 1;
    private background: Phaser.GameObjects.Rectangle;
    private fill: Phaser.GameObjects.Rectangle;
    private trackingSprite: Phaser.GameObjects.Sprite;
    private isDestroyed: boolean = false;

    constructor(props: Props) {
        const { scene, x, y, currentHP, maxHP, trackingSprite } = props;

        this.scene = scene;
        this.trackingSprite = trackingSprite;
        this.maxHP = maxHP;
        this.currentHP = currentHP;

        // Create the background
        this.background = this.scene.add.rectangle(
            x,
            y,
            this.maxFillWidth + 2 * this.padding,
            this.fillHeight + 2 * this.padding,
            0x000000
        ).setOrigin(0, 0).setAlpha(0.3).setDepth(5000);

        // Create the fill bar
        this.fill = this.scene.add.rectangle(
            x + this.padding,
            y + this.padding,
            this.maxFillWidth,
            this.fillHeight,
            0xf5555d
        ).setOrigin(0, 0).setDepth(5001);

        this.setPosition(x, y);
        this.updateHP(currentHP);
    }

    public updateHP(currentHP: number): void {
        if (this.isDestroyed || !this.trackingSprite) {
            this.destroy();
            return;
        }

        this.currentHP = Phaser.Math.Clamp(currentHP, 0, this.maxHP);
        this.fill.setScale(this.currentHP / this.maxHP, 1);
    }

    public setPosition(x: number, y: number): void {
        if (this.isDestroyed) return;

        const offsetX = 32 - (this.maxFillWidth / 2 + this.padding);
        this.background.setPosition(x + offsetX, y);
        this.fill.setPosition(x + offsetX + this.padding, y + this.padding);
    }

    public destroy(): void {
        if (this.isDestroyed) return;

        this.isDestroyed = true;

        this.background?.destroy();
        this.fill?.destroy();

        this.background = null!;
        this.fill = null!;
        this.trackingSprite = null!;
        this.scene = null!;
    }
}
