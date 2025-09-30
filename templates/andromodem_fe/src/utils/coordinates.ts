import type { RelativeCoordinates, TouchMessage } from '@/types/mirroring'
import { MessageType,TouchAction } from '@/types/mirroring'

export interface VideoDisplayArea {
  x: number;
  y: number;
  width: number;
  height: number;
  scaleX: number;
  scaleY: number;
}

export const calculateVideoDisplayArea = (
  containerWidth: number,
  containerHeight: number,
  videoWidth: number,
  videoHeight: number
): VideoDisplayArea => {
  const containerAspect = containerWidth / containerHeight
  const videoAspect = videoWidth / videoHeight
  
  let displayWidth: number
  let displayHeight: number
  let offsetX: number
  let offsetY: number
  
  if (containerAspect > videoAspect) {
    displayHeight = containerHeight
    displayWidth = containerHeight * videoAspect
    offsetX = (containerWidth - displayWidth) / 2
    offsetY = 0
  } else {
    displayWidth = containerWidth
    displayHeight = containerWidth / videoAspect
    offsetX = 0
    offsetY = (containerHeight - displayHeight) / 2
  }
  
  return {
    x: offsetX,
    y: offsetY,
    width: displayWidth,
    height: displayHeight,
    scaleX: displayWidth / videoWidth,
    scaleY: displayHeight / videoHeight
  }
}

export const convertToRelativeCoordinates = (
  clientX: number,
  clientY: number,
  canvasRect: DOMRect,
  videoDisplayArea?: VideoDisplayArea
): RelativeCoordinates | null => {
  const canvasX = clientX - canvasRect.left
  const canvasY = clientY - canvasRect.top
  
  if (!videoDisplayArea) {
    return {
      x: Math.max(0, Math.min(1, canvasX / canvasRect.width)),
      y: Math.max(0, Math.min(1, canvasY / canvasRect.height))
    }
  }
  
  if (canvasX < videoDisplayArea.x || 
      canvasX > videoDisplayArea.x + videoDisplayArea.width ||
      canvasY < videoDisplayArea.y || 
      canvasY > videoDisplayArea.y + videoDisplayArea.height) {
    return null
  }
  
  const videoX = canvasX - videoDisplayArea.x
  const videoY = canvasY - videoDisplayArea.y
  
  return {
    x: Math.max(0, Math.min(1, videoX / videoDisplayArea.width)),
    y: Math.max(0, Math.min(1, videoY / videoDisplayArea.height))
  }
}

export const validateTouchMessage = (message: TouchMessage): boolean => {
  return (
    message.type === MessageType.TOUCH &&
    Object.values(TouchAction).includes(message.action) &&
    message.x >= 0 && message.x <= 1 &&
    message.y >= 0 && message.y <= 1 &&
    message.pointerId >= 0 && message.pointerId <= 10 &&
    message.pressure >= 0 && message.pressure <= 1
  )
}

export const getPressureFromEvent = (event: PointerEvent | MouseEvent | TouchEvent): number => {
  if ('pressure' in event && typeof event.pressure === 'number') {
    return event.pressure
  }
  
  if (event.type.startsWith('mouse')) {
    return event.type === 'mousedown' || event.type === 'mousemove' ? 0.5 : 0
  }
  
  return 0.5
}

export const getPointerIdFromEvent = (event: PointerEvent | MouseEvent | TouchEvent): number => {
  if ('pointerId' in event) {
    return event.pointerId
  }
  
  if ('identifier' in event && typeof event.identifier === 'number') {
    return event.identifier
  }
  
  return 0
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null
  let lastExecTime = 0
  
  return (...args: Parameters<T>) => {
    const currentTime = Date.now()
    
    if (currentTime - lastExecTime > delay) {
      func(...args)
      lastExecTime = currentTime
    } else {
      if (timeoutId) {
        clearTimeout(timeoutId)
      }
      
      timeoutId = setTimeout(() => {
        func(...args)
        lastExecTime = Date.now()
      }, delay - (currentTime - lastExecTime))
    }
  }
}