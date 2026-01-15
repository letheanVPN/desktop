/**
 * Electron ipcRenderer Compatibility Layer
 *
 * Maps Electron's ipcRenderer API to Wails Events and Call system.
 *
 * Electron Concept -> Wails Equivalent:
 * - ipcRenderer.send()    -> Events.Emit() (fire-and-forget)
 * - ipcRenderer.invoke()  -> Call() to bound Go methods (returns Promise)
 * - ipcRenderer.on()      -> Events.On() (subscribe to events)
 * - ipcRenderer.once()    -> Events.Once() (one-time subscription)
 *
 * @example
 * // Electron style:
 * const result = await ipcRenderer.invoke('blockchain:fetchBlock', blockId);
 *
 * // This maps to Wails:
 * const result = await Call.ByName('blockchain.Service.FetchBlockData', blockId);
 */

import { Events, Call } from '@wailsio/runtime';

/**
 * Event object passed to listeners (Electron compatibility)
 */
export interface IpcRendererEvent {
  sender: unknown;
  /** The event name/channel */
  channel: string;
}

export interface IpcMainInvokeEvent extends IpcRendererEvent {
  /** Frame ID - not applicable in Wails, always -1 */
  frameId: number;
}

type IpcListener = (event: IpcRendererEvent, ...args: unknown[]) => void;

// Store for managing event subscriptions (Wails returns cancel functions)
const listenerMap = new Map<string, Map<IpcListener, () => void>>();

/**
 * Electron-compatible ipcRenderer implementation backed by Wails runtime
 */
export const ipcRenderer = {
  /**
   * Send a message to the main process (fire-and-forget).
   * In Wails, this emits an event that Go handlers can listen to.
   *
   * @param channel - The event channel name
   * @param args - Arguments to send
   *
   * @example
   * ipcRenderer.send('user:logout');
   * ipcRenderer.send('analytics:track', { event: 'pageView', page: '/home' });
   */
  send(channel: string, ...args: unknown[]): void {
    Events.Emit({ name: channel, data: args.length === 1 ? args[0] : args });
  },

  /**
   * Send a message and wait for a response (Promise-based).
   * Maps to calling bound Go service methods.
   *
   * IMPORTANT: Channel format determines the Go method called:
   * - 'service:method' -> Attempts to call Service.Method()
   * - 'package.Service.Method' -> Direct Wails binding call
   *
   * @param channel - The channel/method identifier
   * @param args - Arguments to pass to the Go method
   * @returns Promise resolving to the Go method's return value
   *
   * @example
   * // Call blockchain service
   * const block = await ipcRenderer.invoke('blockchain:fetchBlockData', '12345');
   *
   * // Or use direct Wails binding path
   * const block = await ipcRenderer.invoke('blockchain.Service.FetchBlockData', '12345');
   */
  async invoke<T = unknown>(channel: string, ...args: unknown[]): Promise<T> {
    // Convert electron-style 'service:method' to Wails binding path
    const bindingPath = convertChannelToBinding(channel);

    try {
      // Call the bound Go method
      const result = await Call.ByName(bindingPath, ...args);
      return result as T;
    } catch (error) {
      // Wrap in Electron-like error format
      throw new Error(`Error invoking '${channel}': ${error}`);
    }
  },

  /**
   * Subscribe to messages from the main process.
   *
   * @param channel - The event channel to listen on
   * @param listener - Callback function receiving (event, ...args)
   * @returns this (for chaining)
   *
   * @example
   * ipcRenderer.on('mining:hashrate-update', (event, hashrate) => {
   *   console.log('New hashrate:', hashrate);
   * });
   */
  on(channel: string, listener: IpcListener): typeof ipcRenderer {
    const wrappedListener = (data: unknown) => {
      const event: IpcRendererEvent = { sender: null, channel };
      const args = Array.isArray(data) ? data : [data];
      listener(event, ...args);
    };

    // Wails Events.On returns a cancel function
    const cancel = Events.On(channel, wrappedListener);

    // Store the mapping so we can remove it later
    if (!listenerMap.has(channel)) {
      listenerMap.set(channel, new Map());
    }
    listenerMap.get(channel)!.set(listener, cancel);

    return this;
  },

  /**
   * Subscribe to a single message, then auto-unsubscribe.
   *
   * @param channel - The event channel to listen on
   * @param listener - Callback function receiving (event, ...args)
   * @returns this (for chaining)
   *
   * @example
   * ipcRenderer.once('app:ready', (event) => {
   *   console.log('App is ready!');
   * });
   */
  once(channel: string, listener: IpcListener): typeof ipcRenderer {
    const wrappedListener = (data: unknown) => {
      const event: IpcRendererEvent = { sender: null, channel };
      const args = Array.isArray(data) ? data : [data];
      listener(event, ...args);
    };

    Events.Once(channel, wrappedListener);
    return this;
  },

  /**
   * Remove a specific listener from a channel.
   *
   * @param channel - The event channel
   * @param listener - The listener function to remove
   * @returns this (for chaining)
   */
  removeListener(channel: string, listener: IpcListener): typeof ipcRenderer {
    const channelListeners = listenerMap.get(channel);
    if (channelListeners) {
      const cancel = channelListeners.get(listener);
      if (cancel) {
        cancel(); // Call the Wails cancel function
        channelListeners.delete(listener);
      }
    }
    return this;
  },

  /**
   * Remove all listeners for a channel (or all channels if none specified).
   *
   * @param channel - Optional channel to clear; if omitted, clears all
   * @returns this (for chaining)
   */
  removeAllListeners(channel?: string): typeof ipcRenderer {
    if (channel) {
      const channelListeners = listenerMap.get(channel);
      if (channelListeners) {
        channelListeners.forEach((cancel) => cancel());
        listenerMap.delete(channel);
      }
    } else {
      listenerMap.forEach((channelListeners) => {
        channelListeners.forEach((cancel) => cancel());
      });
      listenerMap.clear();
    }
    return this;
  },

  /**
   * Send a synchronous message (NOT RECOMMENDED).
   * Wails doesn't support sync IPC - this throws an error.
   *
   * @deprecated Use invoke() instead for request/response patterns
   * @throws Always throws - sync IPC not supported in Wails
   */
  sendSync(_channel: string, ..._args: unknown[]): never {
    throw new Error(
      'sendSync is not supported in Wails. Use ipcRenderer.invoke() for request/response patterns.'
    );
  },

  /**
   * Post a message to a specific frame (NOT APPLICABLE).
   * Wails doesn't have the same frame concept as Electron.
   *
   * @deprecated Not applicable in Wails architecture
   */
  postMessage(_channel: string, _message: unknown, _transfer?: unknown[]): void {
    console.warn('postMessage is not applicable in Wails. Use send() or invoke() instead.');
  },
};

/**
 * Convert Electron-style channel names to Wails binding paths.
 *
 * Examples:
 * - 'blockchain:fetchBlockData' -> 'github.com/letheanVPN/desktop/services/blockchain.Service.FetchBlockData'
 * - 'config:get' -> 'github.com/letheanVPN/desktop/services/core/config.Service.Get'
 * - 'mining.Service.Start' -> 'github.com/letheanVPN/desktop/services/mining.Service.Start'
 *
 * Add your own mappings in the channelMappings object below!
 */
function convertChannelToBinding(channel: string): string {
  // If it already looks like a binding path, use it directly
  if (channel.includes('.') && !channel.includes(':')) {
    return channel;
  }

  // Known service mappings - ADD YOUR MAPPINGS HERE
  const channelMappings: Record<string, string> = {
    // Blockchain service
    'blockchain:fetchBlockData': 'github.com/letheanVPN/desktop/services/blockchain.Service.FetchBlockData',
    'blockchain:start': 'github.com/letheanVPN/desktop/services/blockchain.Service.Start',
    'blockchain:install': 'github.com/letheanVPN/desktop/services/blockchain.Service.Install',

    // Config service
    'config:get': 'github.com/letheanVPN/desktop/services/core/config.Service.Get',
    'config:isFeatureEnabled': 'github.com/letheanVPN/desktop/services/core/config.Service.IsFeatureEnabled',

    // Display service
    'display:showEnvironmentDialog': 'github.com/letheanVPN/desktop/services/core/display.Service.ShowEnvironmentDialog',
    'display:openWindow': 'github.com/letheanVPN/desktop/services/core/display.Service.OpenWindow',

    // Mining service
    'mining:start': 'github.com/letheanVPN/desktop/services/mining.Service.Start',
    'mining:stop': 'github.com/letheanVPN/desktop/services/mining.Service.Stop',
    'mining:getStats': 'github.com/letheanVPN/desktop/services/mining.Service.GetStats',

    // i18n service
    'i18n:translate': 'github.com/letheanVPN/desktop/services/core/i18n.Service.Translate',

    // Docs service
    'docs:openDocsWindow': 'github.com/letheanVPN/desktop/services/docs.Service.OpenDocsWindow',
  };

  const mapped = channelMappings[channel.toLowerCase()] || channelMappings[channel];
  if (mapped) {
    return mapped;
  }

  // Auto-convert 'service:method' pattern
  // e.g., 'blockchain:fetchBlock' -> 'blockchain.Service.FetchBlock'
  const [service, method] = channel.split(':');
  if (service && method) {
    const pascalMethod = method.charAt(0).toUpperCase() + method.slice(1);
    console.warn(
      `[electron-compat] Auto-converting channel '${channel}'. ` +
      `Consider adding an explicit mapping for better reliability.`
    );
    return `${service}.Service.${pascalMethod}`;
  }

  // Fallback: return as-is and let Wails handle the error
  return channel;
}
