export class EmoticonEmitter {
    private scene: Phaser.Scene;
    private texture: string;
    private frame?: string | number;
    private x: number;
    private y: number;
    // private target?: Phaser.GameObjects.Sprite;
    private timer: Phaser.Time.TimerEvent;

    constructor(
        scene: Phaser.Scene,
        texture: string,
        x: number,
        y: number,
        frame?: string | number,
    ) {
        this.scene = scene;
        this.texture = texture;
        this.frame = frame;
        this.x = x;
        this.y = y;

        this.timer = this.scene.time.addEvent({
            delay: 300,
            callback: this.emit,
            callbackScope: this,
            loop: true
        });
    }

    private emit(): void {
        const offsetX: number = Phaser.Math.Between(-24, 24);
        const offsetY: number = -Phaser.Math.Between(24, 32);

        const sprite = this.scene.add.sprite(this.x + offsetX*0.5, this.y, this.texture, this.frame);

        sprite.setAlpha(1);
        sprite.setDepth(10000);
        sprite.setOrigin(0.5, 0.5);
        sprite.setScale(0.5);

        this.scene.tweens.add({
            targets: sprite,
            y: sprite.y + offsetY,
            duration: 1000,
            ease: 'Back.easeOut',
            onComplete: () => {
                sprite.destroy();
            }
        });

        this.scene.tweens.add({
            targets: sprite,
            x: sprite.x + offsetX,
            scale: 0.4,
            alpha: 0,
            ease: "Linear"
        })
    }

    public setPosition(x: number, y: number) {
        this.x = x;
        this.y = y;
    }

    public stop(): void {
        this.timer.remove(false);
        // this.scene.events.off('update', this.update, this);
    }
}
