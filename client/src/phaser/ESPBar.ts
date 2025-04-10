interface Props {
    scene: Phaser.Scene;
    x: number;
    y: number;
    ecto: number;
    spark: number;
    pulse: number;
    maxESP: number;
    trackingSprite: Phaser.GameObjects.Sprite;
}

// HP bars are scaled on an assumed max 1000 pulse
export class ESPBar {
    private scene: Phaser.Scene;
    private maxESP: number;
    private ecto: number;
    private spark: number;
    private pulse: number;
    private maxFillWidth: number = 32;
    private fillHeight: number = 2;
    private padding: number = 1;

    private ectoBg: Phaser.GameObjects.Rectangle;
    private ectoFill: Phaser.GameObjects.Rectangle;

    private sparkBg: Phaser.GameObjects.Rectangle;
    private sparkFill: Phaser.GameObjects.Rectangle;

    private pulseBg: Phaser.GameObjects.Rectangle;
    private pulseFill: Phaser.GameObjects.Rectangle;

    private trackingSprite: Phaser.GameObjects.Sprite;
    private isDestroyed: boolean = false;

    constructor(props: Props) {
        const { scene, x, y, ecto, spark, pulse, maxESP, trackingSprite } = props;

        this.scene = scene;
        this.trackingSprite = trackingSprite;

        this.maxESP =  maxESP;
        this.ecto = ecto;
        this.spark = spark;
        this.pulse = pulse;

        const ALPHA = 0.5;

        const ectoOffsetY = 0;
        const sparkOffsetY = 2*this.padding + this.fillHeight;
        const pulseOffsetY = 3*this.padding + 2*this.fillHeight;

        // ecto background and fill bar
        this.ectoBg = this.scene.add.rectangle(
            x,
            y + ectoOffsetY,
            this.maxFillWidth + 2 * this.padding,
            this.fillHeight + 1 * this.padding,
            0x000000
        ).setOrigin(0, 0).setAlpha(ALPHA).setDepth(5000);

        this.ectoFill = this.scene.add.rectangle(
            x + this.padding,
            y + sparkOffsetY,
            this.maxFillWidth,
            this.fillHeight,
            0xCA52C9
        ).setOrigin(0, 0).setDepth(5001);

        // spark background and fill bar
        this.sparkBg = this.scene.add.rectangle(
            x,
            y + pulseOffsetY,
            this.maxFillWidth + 2 * this.padding,
            this.fillHeight + 1 * this.padding,
            0x000000
        ).setOrigin(0, 0).setAlpha(ALPHA).setDepth(5000);

        this.sparkFill = this.scene.add.rectangle(
            x + this.padding,
            y - (3*this.padding + 2*this.fillHeight) + this.padding,
            this.maxFillWidth,
            this.fillHeight,
            0x0098DC
        ).setOrigin(0, 0).setDepth(5001);

        // pulse background and fill bar
        this.pulseBg = this.scene.add.rectangle(
            x,
            y - (2*this.padding + 1*this.fillHeight),
            this.maxFillWidth + 2 * this.padding,
            this.fillHeight + 2 * this.padding,
            0x000000
        ).setOrigin(0, 0).setAlpha(ALPHA).setDepth(5000);

        this.pulseFill = this.scene.add.rectangle(
            x + this.padding,
            y - (3*this.padding + 1*this.fillHeight) + this.padding,
            this.maxFillWidth,
            this.fillHeight,
            0xf5555d
        ).setOrigin(0, 0).setDepth(5001);


        this.setPosition(x, y);
        this.updateESP(ecto, spark, pulse);
    }

    public updateESP(ecto: number, spark: number, pulse: number): void {
        if (this.isDestroyed || !this.trackingSprite) {
            this.destroy();
            return;
        }

        this.ecto = Phaser.Math.Clamp(ecto, 0, this.maxESP);
        this.ectoFill.setScale(this.ecto / this.maxESP, 1);

        this.spark = Phaser.Math.Clamp(spark, 0, this.maxESP);
        this.sparkFill.setScale(this.spark / this.maxESP, 1);

        this.pulse = Phaser.Math.Clamp(pulse, 0, this.maxESP);
        this.pulseFill.setScale(this.pulse / this.maxESP, 1);
    }

    public setPosition(x: number, y: number): void {
        if (this.isDestroyed) return;

        const fullHeight = 4*this.padding + 3*this.fillHeight;

        const offsetX = 32 - (this.maxFillWidth / 2 + this.padding);
        var ectoOffsetY = 0 - fullHeight/2;
        var sparkOffsetY = 1*this.padding + 1*this.fillHeight - fullHeight/2;
        var pulseOffsetY = 2*this.padding + 2*this.fillHeight - fullHeight/2;

        this.ectoBg.setPosition(x + offsetX, y + ectoOffsetY);
        this.ectoFill.setPosition(x + offsetX + this.padding, y + ectoOffsetY + this.padding);

        this.sparkBg.setPosition(x + offsetX, y + sparkOffsetY);
        this.sparkFill.setPosition(x + offsetX + this.padding, y + sparkOffsetY + this.padding);

        this.pulseBg.setPosition(x + offsetX, y + pulseOffsetY);
        this.pulseFill.setPosition(x + offsetX + this.padding, y + pulseOffsetY + this.padding);
    }

    public setVisible(visible: boolean) {
        this.ectoBg.setVisible(visible);
        this.ectoFill.setVisible(visible);
        this.sparkBg.setVisible(visible);
        this.sparkFill.setVisible(visible);
        this.pulseBg.setVisible(visible);
        this.pulseFill.setVisible(visible);
    }

    public destroy(): void {
        if (this.isDestroyed) return;

        this.isDestroyed = true;

        this.ectoBg?.destroy();
        this.ectoFill?.destroy();
        this.ectoBg = null!;
        this.ectoFill = null!;

        this.sparkBg?.destroy();
        this.sparkFill?.destroy();
        this.sparkBg = null!;
        this.sparkFill = null!;

        this.pulseBg?.destroy();
        this.pulseFill?.destroy();
        this.pulseBg = null!;
        this.pulseFill = null!;

        this.trackingSprite = null!;
        this.scene = null!;
    }
}
