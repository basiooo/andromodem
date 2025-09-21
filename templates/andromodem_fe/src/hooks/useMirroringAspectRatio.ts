import { useCallback,useEffect } from 'react'

interface UseAspectRatioProps {
  canvasRef: React.RefObject<HTMLCanvasElement>;
  screenWidth?: number;
  screenHeight?: number;
}

export const useAspectRatio = ({ 
  canvasRef, 
  screenWidth, 
  screenHeight 
}: UseAspectRatioProps) => {
  
  const updateAspectRatio = useCallback(() => {
    if (!canvasRef.current || !screenWidth || !screenHeight) return
    
    const canvas = canvasRef.current
    const container = canvas.parentElement
    
    if (container) {
      const aspectRatio = screenWidth / screenHeight
      container.style.aspectRatio = aspectRatio.toString()
    }
  }, [canvasRef, screenWidth, screenHeight])

  useEffect(() => {
    updateAspectRatio()
  }, [updateAspectRatio])

  return { updateAspectRatio }
}