import Phaser from "phaser";
import { BaseEntity } from "./entities/BaseEntity"; // Updated import path

export class SelectionManager {
    private scene: Phaser.Scene;
    private selectedEntity: BaseEntity | null = null;
    private isPanning: boolean = false;
    private holdTimer: Phaser.Time.TimerEvent | null = null;
    private readonly holdThreshold: number = 200; // ms to distinguish click vs hold
    private startPointer: Phaser.Input.Pointer | null = null;

    constructor(scene: Phaser.Scene) {
        this.scene = scene;

        this.scene.input.on("gameobjectdown", this.handleObjectDown, this);
        this.scene.input.on("gameobjectup", this.handleObjectUp, this);

        this.scene.input.on("pointerdown", this.handlePointerDown, this);
        this.scene.input.on("pointermove", this.handlePointerMove, this);
        this.scene.input.on("pointerup", this.handlePointerUp, this);
    }

    private handleObjectDown(pointer: Phaser.Input.Pointer, gameObject: Phaser.GameObjects.GameObject) {
        this.startPointer = pointer;
        this.isPanning = false;

        this.holdTimer = this.scene.time.delayedCall(this.holdThreshold, () => {
            if (this.startPointer?.isDown) {
                this.isPanning = true;
            }
        }, [], this);
    }

    private handleObjectUp(pointer: Phaser.Input.Pointer, gameObject: Phaser.GameObjects.GameObject) {
        if (!this.startPointer || this.isPanning) return;

        const entity = gameObject.getData("entity") as BaseEntity;
        if (entity instanceof BaseEntity && pointer.time - this.startPointer.downTime < this.holdThreshold) {
            this.setSelected(entity);
        }
    }

    private handlePointerDown(pointer: Phaser.Input.Pointer) {
        this.startPointer = pointer;
        this.isPanning = false;

        this.holdTimer = this.scene.time.delayedCall(this.holdThreshold, () => {
            if (this.startPointer?.isDown) {
                this.isPanning = true;
            }
        }, [], this);
    }

    private handlePointerMove(pointer: Phaser.Input.Pointer) {
        if (!this.isPanning || !this.startPointer || pointer !== this.startPointer) return;

        const dx = pointer.x - pointer.prevPosition.x;
        const dy = pointer.y - pointer.prevPosition.y;
        this.scene.cameras.main.scrollX -= dx;
        this.scene.cameras.main.scrollY -= dy;
    }

    private handlePointerUp(pointer: Phaser.Input.Pointer) {
        if (!this.startPointer || pointer !== this.startPointer) return;

        if (!this.isPanning && pointer.time - this.startPointer.downTime < this.holdThreshold) {
            const hitObjects = this.scene.input.manager.hitTest(pointer, this.scene.children.list, this.scene.cameras.main);
            if (hitObjects.length === 0 && this.selectedEntity) {
                this.setSelected(null);
            }
        }

        this.isPanning = false;
        this.startPointer = null;
        if (this.holdTimer) {
            this.holdTimer.remove(false);
            this.holdTimer = null;
        }
    }

    private setSelected(entity: BaseEntity | null) {
        if (this.selectedEntity === entity) return;

        if (this.selectedEntity) {
            this.selectedEntity.setSelected(false);
        }

        this.selectedEntity = entity;
        if (entity) {
            entity.setSelected(true);
        }

        const eventData = entity ? { id: entity.id, type: entity.type, data: entity.data } : null;
        window.dispatchEvent(new CustomEvent("entitySelection", { detail: eventData }));
    }
}