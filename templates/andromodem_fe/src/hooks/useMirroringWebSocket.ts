import { useEffect, useRef, useState, useCallback } from 'react'

import { config } from '@/config'
import { DeviceState } from '@/types/device'
import type { 
  ConnectedMessage, 
  ConnectionStateValue, 
  KeyMessage,
  SetupMirroring,
  TouchMessage,
  UseMirroringWebSocketOptions, 
  UseMirroringWebSocketReturn, 
  WebSocketMessage} from '@/types/mirroring'
import { ConnectionState, KeyCommand,MessageType } from '@/types/mirroring'
import { validateTouchMessage } from '@/utils/coordinates'

const PING_INTERVAL = 5000 

export const useMirroringWebSocket = (options: UseMirroringWebSocketOptions): UseMirroringWebSocketReturn =>  {
  const { device, onConnected, onError, onVideoFrame } = options
  
  const [connectionState, setConnectionState] = useState<ConnectionStateValue>(ConnectionState.DISCONNECTED)
  const [error, setError] = useState<string | null>(null)
  const [screenDimensions, setScreenDimensions] = useState<{ width: number; height: number } | null>(null)
  
  const wsRef = useRef<WebSocket | null>(null)
  const pingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const isManualDisconnectRef = useRef(false)
  
  const isConnected = connectionState === ConnectionState.CONNECTED
  const isConnecting = connectionState === ConnectionState.CONNECTING
  
  const clearTimeouts = () => {
    if (pingIntervalRef.current) {
      clearInterval(pingIntervalRef.current)
      pingIntervalRef.current = null
    }
  }
  
  const setupPingInterval = () => {
    clearTimeouts()
    pingIntervalRef.current = setInterval(() => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'ping' }))
      }
    }, PING_INTERVAL)
  }
  
  const handleMessage = (event: MessageEvent) => {
    if (event.data instanceof ArrayBuffer) {
      onVideoFrame?.(event.data)
      return
    }
    
    try {
      const message: WebSocketMessage = JSON.parse(event.data)
      
      switch (message.type) {
        case MessageType.CONNECTED:
          const connectedMsg = message as ConnectedMessage
          setConnectionState(ConnectionState.CONNECTED)
          setError(null)
          setScreenDimensions({ 
            width: connectedMsg.width, 
            height: connectedMsg.height 
          })
          setupPingInterval()
          onConnected?.()
          break
          
        case MessageType.ERROR:
          setError(message.message)
          setConnectionState(ConnectionState.ERROR)
          onError?.(message.message)
          break
          
        default:
          console.warn('Unknown message type:', message)
      }
    } catch (err) {
      console.error('Failed to parse WebSocket message:', err)
      setError('Failed to parse server message')
    }
  }
  
  const handleError = (event: Event) => {
    console.error('WebSocket error:', event)
    setError('WebSocket connection error')
    setConnectionState(ConnectionState.ERROR)
    onError?.('WebSocket connection error')
  }
  
  const handleClose = (event: CloseEvent) => {
    console.log('WebSocket closed:', event.code, event.reason)
    setConnectionState(ConnectionState.DISCONNECTED)
    clearTimeouts()
  }
  
  const connect = useCallback((setup: SetupMirroring) => {
    if (
      device.state !== DeviceState.ONLINE ||
      wsRef.current?.readyState === WebSocket.OPEN ||
      wsRef.current?.readyState === WebSocket.CONNECTING
    ) {
      return
    }

    isManualDisconnectRef.current = false
    setConnectionState(ConnectionState.CONNECTING)
    setError(null)

    try {
      const wsUrl = `${config.BASE_URL_WS}/ws/devices/${device.serial}/mirroring`
      wsRef.current = new WebSocket(wsUrl)
      wsRef.current.binaryType = 'arraybuffer'

      wsRef.current.onopen = () => {
        console.log('WebSocket connected')
        wsRef.current?.send(JSON.stringify(setup))
      }

      wsRef.current.onmessage = handleMessage
      wsRef.current.onerror = handleError
      wsRef.current.onclose = handleClose
    } catch (err) {
      console.error('Failed to create WebSocket:', err)
      setError('Failed to create WebSocket connection')
      setConnectionState(ConnectionState.ERROR)
    }
  }, [device.serial, device.state])
  
  const disconnect = useCallback(() => {
    isManualDisconnectRef.current = true
    clearTimeouts()

    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
    setConnectionState(ConnectionState.DISCONNECTED)
    setError(null)
    setScreenDimensions(null)
  }, [])
  
  const sendTouchEvent = (event: TouchMessage) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket not connected, cannot send touch event')
      return
    }
    
    if (!validateTouchMessage(event)) {
      console.warn('Invalid touch message:', event)
      return
    }
    
    try {
      wsRef.current.send(JSON.stringify(event))
    } catch (err) {
      console.error('Failed to send touch event:', err)
      setError('Failed to send touch event')
    }
  }

  const sendKeyEvent = (event: KeyMessage) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket not connected, cannot send key event')
      return
    }
    
    const allowedKeys = Object.values(KeyCommand)
    if (!allowedKeys.includes(event.key)) {
      console.warn('Invalid key:', event.key)
      return
    }
    
    try {
      wsRef.current.send(JSON.stringify(event))
      console.log('Key event sent:', event.key)
    } catch (err) {
      console.error('Failed to send key event:', err)
      setError('Failed to send key event')
    }
  }
  
  useEffect(() => {
    return () => {
      disconnect()
    }
  }, [])
  
  return {
    isConnected,
    isConnecting,
    error,
    sendTouchEvent,
    sendKeyEvent,
    connect,
    disconnect,
    screenDimensions
  }
}