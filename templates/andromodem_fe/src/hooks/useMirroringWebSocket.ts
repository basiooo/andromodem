import { useCallback,useEffect, useRef, useState } from 'react'

import { config } from '@/config'
import { DeviceState } from '@/types/device'
import type { 
  ConnectedMessage, 
  KeyMessage,
  TouchMessage,
  UseMirroringWebSocketOptions, 
  UseMirroringWebSocketReturn, 
  WebSocketMessage} from '@/types/mirroring'
import { ConnectionState, KeyCommand,MessageType } from '@/types/mirroring'
import { validateTouchMessage } from '@/utils/coordinates'

const PING_INTERVAL = 5000 

export const useMirroringWebSocket = (options: UseMirroringWebSocketOptions): UseMirroringWebSocketReturn =>  {
  const { device, onConnected, onError, onVideoFrame } = options
  
  const [connectionState, setConnectionState] = useState<ConnectionState>(ConnectionState.DISCONNECTED)
  const [error, setError] = useState<string | null>(null)
  const [screenDimensions, setScreenDimensions] = useState<{ width: number; height: number } | null>(null)
  
  const wsRef = useRef<WebSocket | null>(null)
  const pingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const isManualDisconnectRef = useRef(false)
  
  const isConnected = connectionState === ConnectionState.CONNECTED
  const isConnecting = connectionState === ConnectionState.CONNECTING
  
  const clearTimeouts = useCallback(() => {
    if (pingIntervalRef.current) {
      clearInterval(pingIntervalRef.current)
      pingIntervalRef.current = null
    }
  }, [])
  
  const setupPingInterval = useCallback(() => {
    clearTimeouts()
    pingIntervalRef.current = setInterval(() => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'ping' }))
      }
    }, PING_INTERVAL)
  }, [clearTimeouts])
  
  const handleMessage = useCallback((event: MessageEvent) => {
    if (event.data instanceof ArrayBuffer) {
      onVideoFrame?.(event.data)
      return
    }
    
    try {
      const message: WebSocketMessage = JSON.parse(event.data)
      
      switch (message.type) {
        case MessageType.CONNECTED:
          // eslint-disable-next-line no-case-declarations
          const connectedMsg = message as ConnectedMessage
          setConnectionState(ConnectionState.CONNECTED)
          setError(null)
          setScreenDimensions({ 
            width: connectedMsg.width, 
            height: connectedMsg.height 
          })
          setupPingInterval()
          onConnected?.(connectedMsg)
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
  }, [onConnected, onError, onVideoFrame, setupPingInterval])
  
  const handleError = useCallback((event: Event) => {
    console.error('WebSocket error:', event)
    setError('WebSocket connection error')
    setConnectionState(ConnectionState.ERROR)
    onError?.('WebSocket connection error')
  }, [onError])
  
  const handleClose = useCallback((event: CloseEvent) => {
    console.log('WebSocket closed:', event.code, event.reason)
    setConnectionState(ConnectionState.DISCONNECTED)
    clearTimeouts()
  }, [clearTimeouts])
  
  const connect = useCallback(() => {
    if (device.state !== DeviceState.ONLINE || wsRef.current?.readyState === WebSocket.OPEN) {
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
      }
      
      wsRef.current.onmessage = handleMessage
      wsRef.current.onerror = handleError
      wsRef.current.onclose = handleClose
      
    } catch (err) {
      console.error('Failed to create WebSocket:', err)
      setError('Failed to create WebSocket connection')
      setConnectionState(ConnectionState.ERROR)
    }
  }, [device, handleMessage, handleError, handleClose])
  
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
  }, [clearTimeouts])
  
  const sendTouchEvent = useCallback((event: TouchMessage) => {
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
  }, [])

  const sendKeyEvent = useCallback((event: KeyMessage) => {
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
  }, [])
  
  useEffect(() => {
    return () => {
      disconnect()
    }
  }, [disconnect])
  
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