type EventMap = {
    panToGotchi: { gotchiId: string };
    entitySelection: { detail: any };
  };
  
  type Listener<T> = (data: T) => void;
  
  export class EventBus {
    private listeners: {
      [K in keyof EventMap]: Set<Listener<EventMap[K]>>;
    } = {
      panToGotchi: new Set(),
      entitySelection: new Set(),
    };
  
    on<K extends keyof EventMap>(event: K, listener: Listener<EventMap[K]>) {
      this.listeners[event].add(listener);
    }
  
    off<K extends keyof EventMap>(event: K, listener: Listener<EventMap[K]>) {
      this.listeners[event].delete(listener);
    }
  
    emit<K extends keyof EventMap>(event: K, data: EventMap[K]) {
      // Check if there are listeners for the event
      if (this.listeners[event].size === 0) {
        console.log(`[EventBus] No listeners for event: ${event}`);
      }
  
      // Emit the event to all listeners
      for (const listener of this.listeners[event]) {
        listener(data);
      }
    }
  }
  
  export const eventBus = new EventBus();
  