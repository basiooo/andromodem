import { useCallback, useEffect, useRef, useState } from 'react'

import type {
  TouchMessage,
  UseMirroringTouchOptions,
  UseMirroringTouchReturn
} from '@/types/mirroring'
import { MessageType,TouchAction } from '@/types/mirroring'
import {
  calculateVideoDisplayArea,
  convertToRelativeCoordinates,
  throttle,
  type VideoDisplayArea
} from '@/utils/coordinates'

const TOUCH_THROTTLE_MS = 16

export const useMonitoringTouch = (options: UseMirroringTouchOptions): UseMirroringTouchReturn => {
  const { canvasRef, screenWidth, screenHeight, onTouchEvent, enabled = true } = options

  const [isActive, setIsActive] = useState(false)
  const [activePointers] = useState(new Map<number, { x: number; y: number }>())
  const [videoDisplayArea, setVideoDisplayArea] = useState<VideoDisplayArea | undefined>(undefined)
  const isMouseDownRef = useRef(false)

  const throttledSendTouchEvent = useRef(
    throttle((event: TouchMessage) => {
      onTouchEvent(event)
    }, TOUCH_THROTTLE_MS)
  ).current

  useEffect(() => {
    const updateVideoArea = () => {
      const canvas = canvasRef.current
      if (!canvas) return

      const container = canvas.parentElement
      if (!container) return

      const rect = container.getBoundingClientRect()
      const area = calculateVideoDisplayArea(
        rect.width,
        rect.height,
        screenWidth,
        screenHeight
      )
      setVideoDisplayArea(area)
    }

    updateVideoArea()

    const canvas = canvasRef.current
    if (!canvas || !canvas.parentElement) return

    const resizeObserver = new ResizeObserver(updateVideoArea)
    resizeObserver.observe(canvas.parentElement)

    return () => {
      resizeObserver.disconnect()
    }
  }, [canvasRef, screenWidth, screenHeight])

  const createTouchMessage = useCallback((
    action: TouchAction,
    clientX: number,
    clientY: number,
    pointerId: number,
    pressure: number
  ): TouchMessage | null => {
    if (!canvasRef.current) return null

    const container = canvasRef.current.parentElement
    if (!container) return null

    const rect = container.getBoundingClientRect()
    const coords = convertToRelativeCoordinates(clientX, clientY, rect, videoDisplayArea)

    if (!coords) return null

    return {
      type: MessageType.TOUCH,
      action,
      x: coords.x,
      y: coords.y,
      pointerId,
      pressure
    }
  }, [canvasRef, videoDisplayArea])

  const updateActivePointers = useCallback((pointerId: number, x: number, y: number, remove = false) => {
    if (remove) {
      activePointers.delete(pointerId)
    } else {
      activePointers.set(pointerId, { x, y })
    }
    setIsActive(activePointers.size > 0)
  }, [activePointers])

  const handleTouchStart = useCallback((event: TouchEvent) => {
    if (!enabled || !canvasRef.current) return

    event.preventDefault()
    setIsActive(true)

    for (let i = 0; i < event.changedTouches.length; i++) {
      const touch = event.changedTouches[i]
      const touchMessage = createTouchMessage(TouchAction.DOWN, touch.clientX, touch.clientY, touch.identifier, 0.5)

      if (touchMessage) {
        updateActivePointers(touch.identifier, touchMessage.x, touchMessage.y)
        onTouchEvent(touchMessage)
      }
    }
  }, [enabled, canvasRef, createTouchMessage, updateActivePointers, onTouchEvent])

  const handleTouchMove = useCallback((event: TouchEvent) => {
    if (!enabled || !canvasRef.current) return

    event.preventDefault()

    for (let i = 0; i < event.changedTouches.length; i++) {
      const touch = event.changedTouches[i]
      const touchMessage = createTouchMessage(TouchAction.MOVE, touch.clientX, touch.clientY, touch.identifier, 0.5)

      if (touchMessage) {
        updateActivePointers(touch.identifier, touchMessage.x, touchMessage.y)
        throttledSendTouchEvent(touchMessage)
      } else {
        const cancelMessage = createTouchMessage(TouchAction.CANCEL, touch.clientX, touch.clientY, touch.identifier, 0)
        if (cancelMessage) {
          onTouchEvent(cancelMessage)
        }
      }
    }
  }, [enabled, canvasRef, createTouchMessage, updateActivePointers, throttledSendTouchEvent, activePointers, onTouchEvent])

  const handleTouchEnd = useCallback((event: TouchEvent) => {
    if (!enabled || !canvasRef.current) return

    event.preventDefault()

    for (let i = 0; i < event.changedTouches.length; i++) {
      const touch = event.changedTouches[i]
      const touchMessage = createTouchMessage(TouchAction.UP, touch.clientX, touch.clientY, touch.identifier, 0)

      if (touchMessage) {
        updateActivePointers(touch.identifier, touchMessage.x, touchMessage.y, true)
        onTouchEvent(touchMessage)
      }
    }

    if (activePointers.size === 0) {
      setIsActive(false)
    }
  }, [enabled, canvasRef, createTouchMessage, updateActivePointers, onTouchEvent, activePointers])

  const handleTouchCancel = useCallback((event: TouchEvent) => {
    if (!enabled || !canvasRef.current) return

    event.preventDefault()

    for (let i = 0; i < event.changedTouches.length; i++) {
      const touch = event.changedTouches[i]
      const touchMessage = createTouchMessage(TouchAction.CANCEL, touch.clientX, touch.clientY, touch.identifier, 0)

      if (touchMessage) {
        updateActivePointers(touch.identifier, touchMessage.x, touchMessage.y, true)
        onTouchEvent(touchMessage)
      }
    }
  }, [enabled, canvasRef, createTouchMessage, updateActivePointers, onTouchEvent])

  const handleMouseDown = useCallback((event: MouseEvent) => {
    if (!enabled || !canvasRef.current) return

    event.preventDefault()
    isMouseDownRef.current = true
    const touchMessage = createTouchMessage(TouchAction.DOWN, event.clientX, event.clientY, 0, 0.5)

    if (touchMessage) {
      setIsActive(true)
      updateActivePointers(0, touchMessage.x, touchMessage.y)
      onTouchEvent(touchMessage)
    }
  }, [enabled, canvasRef, createTouchMessage, updateActivePointers, onTouchEvent])

  const handleMouseMove = useCallback((event: MouseEvent) => {
    if (!enabled || !canvasRef.current || !isMouseDownRef.current) return

    event.preventDefault()
    const touchMessage = createTouchMessage(TouchAction.MOVE, event.clientX, event.clientY, 0, 0.5)

    if (touchMessage) {
      updateActivePointers(0, touchMessage.x, touchMessage.y)
      throttledSendTouchEvent(touchMessage)
    } else {
      const cancelMessage = createTouchMessage(TouchAction.CANCEL, event.clientX, event.clientY, 0, 0)
      if (cancelMessage) {
        onTouchEvent(cancelMessage)
      }
    }
  }, [enabled, canvasRef, createTouchMessage, updateActivePointers, throttledSendTouchEvent, onTouchEvent])

  const handleMouseUp = useCallback((event: MouseEvent) => {
    if (!enabled || !canvasRef.current || !isMouseDownRef.current) return

    event.preventDefault()
    isMouseDownRef.current = false
    const touchMessage = createTouchMessage(TouchAction.UP, event.clientX, event.clientY, 0, 0)

    if (touchMessage) {
      updateActivePointers(0, touchMessage.x, touchMessage.y, true)
      onTouchEvent(touchMessage)
    }

    setIsActive(false)
  }, [enabled, canvasRef, createTouchMessage, updateActivePointers, onTouchEvent])

  const handleMouseLeave = useCallback((event: MouseEvent) => {
    if (!enabled || !canvasRef.current || !isMouseDownRef.current) return

    event.preventDefault()
    isMouseDownRef.current = false
    const touchMessage = createTouchMessage(TouchAction.CANCEL, event.clientX, event.clientY, 0, 0)

    if (touchMessage) {
      updateActivePointers(0, touchMessage.x, touchMessage.y, true)
      onTouchEvent(touchMessage)
    }
  }, [enabled, canvasRef, createTouchMessage, updateActivePointers, onTouchEvent])

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas || !enabled) return

    canvas.addEventListener('touchstart', handleTouchStart, { passive: false })
    canvas.addEventListener('touchmove', handleTouchMove, { passive: false })
    canvas.addEventListener('touchend', handleTouchEnd, { passive: false })
    canvas.addEventListener('touchcancel', handleTouchCancel, { passive: false })
    canvas.addEventListener('mousedown', handleMouseDown)
    canvas.addEventListener('mousemove', handleMouseMove)
    canvas.addEventListener('mouseup', handleMouseUp)
    canvas.addEventListener('mouseleave', handleMouseLeave)

    return () => {
      canvas.removeEventListener('touchstart', handleTouchStart)
      canvas.removeEventListener('touchmove', handleTouchMove)
      canvas.removeEventListener('touchend', handleTouchEnd)
      canvas.removeEventListener('touchcancel', handleTouchCancel)
      canvas.removeEventListener('mousedown', handleMouseDown)
      canvas.removeEventListener('mousemove', handleMouseMove)
      canvas.removeEventListener('mouseup', handleMouseUp)
      canvas.removeEventListener('mouseleave', handleMouseLeave)
    }
  }, [
    enabled,
    canvasRef,
    handleTouchStart,
    handleTouchMove,
    handleTouchEnd,
    handleTouchCancel,
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
    handleMouseLeave
  ])

  return {
    isActive,
    activePointers
  }
}