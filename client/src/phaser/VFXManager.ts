export class VFXManager {
    private static instance: VFXManager | null = null;
    private scene: Phaser.Scene;
    private lickExplosionPool: Phaser.GameObjects.Sprite[] = [];
    private poolSize: number = 10;
    private activeLickExplosions: Set<Phaser.GameObjects.Sprite> = new Set();

    constructor(scene: Phaser.Scene) {
        this.scene = scene;

        // Pre-populate lick explosion pool
        for (let i = 0; i < this.poolSize; i++) {
            const sprite = this.createLickExplosionSprite();
            sprite.setVisible(false);
            this.lickExplosionPool.push(sprite);
        }

        this.createLickExplosionAnim();
    }

        // This static method ensures there is only one instance of VFXManager
        public static getInstance(scene: Phaser.Scene): VFXManager {
            if (!this.instance) {
                this.instance = new VFXManager(scene);
            }
            return this.instance;
        }

    public static preload(scene: Phaser.Scene) {
        scene.load.spritesheet("lickquidator_explosion", "assets/vfx/lick_death_spritesheet.png", {
            frameWidth: 64,
            frameHeight: 64,
            margin: 0,
            spacing: 0
        });
        
    }

    private createLickExplosionSprite(): Phaser.GameObjects.Sprite {
        const sprite = this.scene.add.sprite(0, 0, "lickquidator_explosion").setDepth(9999);
        sprite.setOrigin(0, 0);
        sprite.setVisible(false);
        return sprite;
    }

    private createLickExplosionAnim(): void {
        if (!this.scene.anims.exists("lickquidator_explosion")) {
            this.scene.anims.create({
                key: "lickquidator_explosion",
                frames: this.scene.anims.generateFrameNumbers("lickquidator_explosion", {}),
                frameRate: 16,
                repeat: 0,
                hideOnComplete: true
            });
        }
    }

    public playLickExplosion(x: number, y: number): void {
        const explosion = this.getLickExplosionSprite();
        explosion.setPosition(x, y);
        explosion.setVisible(true);
        explosion.play("lickquidator_explosion");

        explosion.once(Phaser.Animations.Events.ANIMATION_COMPLETE, () => {
            explosion.setVisible(false);
            this.releaseLickExplosionSprite(explosion);
        });
    }

    private getLickExplosionSprite(): Phaser.GameObjects.Sprite {
        if (this.lickExplosionPool.length > 0) {
            const sprite = this.lickExplosionPool.pop()!;
            this.activeLickExplosions.add(sprite);
            return sprite;
        }

        // Pool exhausted, create new sprite
        const newSprite = this.createLickExplosionSprite();
        this.activeLickExplosions.add(newSprite);
        return newSprite;
    }

    private releaseLickExplosionSprite(sprite: Phaser.GameObjects.Sprite): void {
        this.activeLickExplosions.delete(sprite);
        this.lickExplosionPool.push(sprite);
    }

    public destroy(): void {
        this.lickExplosionPool.forEach(sprite => sprite.destroy());
        this.activeLickExplosions.forEach(sprite => sprite.destroy());
        this.lickExplosionPool = [];
        this.activeLickExplosions.clear();
    }
}
