// src/utils/EventBus.ts
type EventMap = {
    panToGotchi: { gotchiId: string };
    // Add other events here as needed
};

type Listener<T> = (data: T) => void;

export class EventBus {
    private listeners: { [K in keyof EventMap]?: Listener<EventMap[K]>[] } = {};

    on<K extends keyof EventMap>(event: K, listener: Listener<EventMap[K]>) {
        if (!this.listeners[event]) {
            this.listeners[event] = [];
        }
        this.listeners[event]!.push(listener);
    }

    off<K extends keyof EventMap>(event: K, listener: Listener<EventMap[K]>) {
        if (!this.listeners[event]) return;
        this.listeners[event] = this.listeners[event]!.filter(l => l !== listener);
    }

    emit<K extends keyof EventMap>(event: K, data: EventMap[K]) {
        if (!this.listeners[event]) return;
        this.listeners[event]!.forEach(listener => listener(data));
    }
}

export const eventBus = new EventBus(); // Singleton instance