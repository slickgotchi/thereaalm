export class EmoticonEmitter {
    private scene: Phaser.Scene;
    private x: number;
    private y: number;
    private isEmitting: boolean = false;

    constructor(scene: Phaser.Scene, x: number, y: number) {
        this.scene = scene;
        this.x = x;
        this.y = y;
    }

    public static preload(scene: Phaser.Scene): void {
        scene.load.spritesheet(
            "emoticons",
            "assets/emoticons/emoticons_48px.png",
            {frameWidth: 48, frameHeight: 48, margin: 2, spacing: 4}
        );
        scene.load.spritesheet(
            "actionicons",
            "assets/spritesheets/actionicons_spritesheet.png",
            {frameWidth: 48, frameHeight: 48}
        );
    }

    private getTextureAndFrame(emoticonStr: string): { texture: string; frame?: string | number } {
        switch (emoticonStr) {
            case "attack": return { texture: 'actionicons', frame: 0 };
            case "forage": return { texture: 'actionicons', frame: 1 };
            case "chop": return { texture: 'actionicons', frame: 2 };
            case "mine": return { texture: 'actionicons', frame: 3 };
            case "harvest": return { texture: 'actionicons', frame: 4 };
            case "flee": return { texture: 'actionicons', frame: 10 };
            case "roam": return { texture: 'actionicons', frame: 12 };
            case "sell": return { texture: 'actionicons', frame: 8 };
            case "buy": return { texture: 'actionicons', frame: 9 };
            case "rest": return { texture: 'actionicons', frame: 11 };
            case "repair": return { texture: 'actionicons', frame: 7 };
            case "build": return { texture: 'actionicons', frame: 6 };
            default: {
                console.log(`No emoticon for '${emoticonStr}'`);
                return { texture: 'icons', frame: 'default' };
            } 
        }
    }

    public emit(emoticon: string, duration_ms: number): void {
        if (this.isEmitting) return;
        this.isEmitting = true;

        const { texture, frame } = this.getTextureAndFrame(emoticon);
        const startTime = this.scene.time.now;

        const emitAction = () => {
            if (duration_ms >= 0 && this.scene.time.now - startTime > duration_ms) {
                this.isEmitting = false;
                return;
            }

            const offsetX: number = Phaser.Math.Between(-24, 24);
            const offsetY: number = -Phaser.Math.Between(24, 32);

            const sprite = this.scene.add.sprite(this.x + offsetX * 0.5, this.y, texture, frame);
            sprite.setAlpha(1).setDepth(10000).setOrigin(0.5, 0.5).setScale(0.5);

            this.scene.tweens.add({
                targets: sprite,
                y: sprite.y + offsetY,
                duration: 1000,
                ease: 'Back.easeOut',
                onComplete: () => sprite.destroy()
            });

            this.scene.tweens.add({
                targets: sprite,
                x: sprite.x + offsetX,
                scale: 0.4,
                ease: "Linear"
            });

            this.scene.tweens.add({
                targets: sprite,
                alpha: 0,
                ease: "Quint.easeIn"
            });

            if (duration_ms < 0 || this.scene.time.now - startTime <= duration_ms) {
                this.scene.time.delayedCall(300, emitAction);
            } else {
                this.isEmitting = false;
            }
        };

        emitAction();
    }

    public setPosition(x: number, y: number) {
        this.x = x;
        this.y = y;
    }

    public destroy() {
        
    }
}
