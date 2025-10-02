import type { Device } from "./device"

export const TouchAction = {
  DOWN: "down",
  UP: "up",
  MOVE: "move",
  CANCEL: "cancel"
} as const

export type TouchActionValue = typeof TouchAction[keyof typeof TouchAction];

export const KeyCommand = {
  BACK: "back",
  HOME: "home",
  RECENT: "recent",
  POWER: "power"
} as const

export type KeyCommandValue = typeof KeyCommand[keyof typeof KeyCommand];

export const MessageType = {
  TOUCH: "touch",
  KEY: "key",
  PING: "ping",
  SETUP: "setup",
  CONNECTED: "connected",
  ERROR: "error"
} as const

export type MessageTypeValue = typeof MessageType[keyof typeof MessageType];

export const ConnectionState = {
  DISCONNECTED: "disconnected",
  CONNECTING: "connecting",
  CONNECTED: "connected",
  ERROR: "error"
} as const

export type ConnectionStateValue = typeof ConnectionState[keyof typeof ConnectionState];

export interface TouchMessage {
  type: typeof MessageType.TOUCH;
  action: TouchActionValue;
  x: number; 
  y: number;
  pointerId: number;
  pressure: number;
}

export interface KeyMessage {
  type: typeof MessageType.KEY;
  key: KeyCommandValue;
}

export interface ConnectedMessage {
  type: typeof MessageType.CONNECTED;
  message: "Mirroring stream connected";
  serial: string;
  width: number;
  height: number;
}

export interface ErrorMessage {
  type: typeof MessageType.ERROR;
  message: string;
}

export type WebSocketMessage = ConnectedMessage | ErrorMessage;

export interface SetupMirroring {
  type: typeof MessageType.SETUP;
  fps: FPSValue;
  resolution: ScreenResolutionValue;
  bitrate: BitRateValue;
}

export interface UseMirroringWebSocketOptions {
  device: Device;
  onConnected: () => void;
  onError: (error: string) => void;
  onVideoFrame: (frame: ArrayBuffer) => void;
}

export interface UseMirroringWebSocketReturn {
  isConnected: boolean;
  isConnecting: boolean;
  error: string | null;
  sendTouchEvent: (event: TouchMessage) => void;
  sendKeyEvent: (event: KeyMessage) => void;
  connect: (options: SetupMirroring) => void;
  disconnect: () => void;
  screenDimensions: { width: number; height: number } | null;
}

export interface UseMirroringTouchOptions {
  canvasRef: React.RefObject<HTMLCanvasElement>;
  screenWidth: number;
  screenHeight: number;
  onTouchEvent: (event: TouchMessage) => void;
  enabled?: boolean;
}

export interface UseMirroringTouchReturn {
  isActive: boolean;
  activePointers: Map<number, { x: number; y: number }>;
}

export interface MirroringCanvasProps {
  device: Device;
  onConnectionChange?: (connected: boolean) => void;
  onError?: (error: string) => void;
}

export interface RelativeCoordinates {
  x: number;
  y: number;
}

export const ScreenResolution = {
  360: 360,
  480: 480,
  720: 720,
  1080: 1080,
} as const
export type ScreenResolutionValue = typeof ScreenResolution[keyof typeof ScreenResolution];

export const BitRate = {
  1: 1000000,
  2: 2000000,
  3: 3000000,
  4: 4000000,
  5: 5000000,
  6: 6000000,
  7: 7000000,
  8: 8000000,
} as const
export type BitRateValue = typeof BitRate[keyof typeof BitRate];

export const FPS = {
  30: 30,
  60: 60,
} as const
export type FPSValue = typeof FPS[keyof typeof FPS];
