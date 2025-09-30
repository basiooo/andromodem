import JMuxer from 'jmuxer'
import { useCallback, useEffect, useRef, useState } from 'react'
import { FaChromecast } from "react-icons/fa"
import { toast } from 'react-toastify'

import MirroringNavigation from '@/components/Mirroring/MirroringNavigation'
import MirroringTool from '@/components/Mirroring/MirroringTool'
import { useAspectRatio } from '@/hooks/useMirroringAspectRatio'
import { useMirroringFullscreen } from '@/hooks/useMirroringFullscreen'
import { useMonitoringTouch } from '@/hooks/useMirroringTouch'
import { useMirroringWebSocket } from '@/hooks/useMirroringWebSocket'
import type { MirroringCanvasProps } from '@/types/mirroring'
import { KeyCommand, MessageType } from '@/types/mirroring'

const MirroringCanvas: React.FC<MirroringCanvasProps> = ({
    device
}) => {
    const canvasRef = useRef<HTMLCanvasElement>(null)
    const videoRef = useRef<HTMLVideoElement>(null)
    const jmuxerRef = useRef<JMuxer | null>(null)
    const mainContainerRef = useRef<HTMLDivElement>(null)
    const [isVideoReady, setIsVideoReady] = useState(false)
    const [isDisconnecting, setIsDisconnecting] = useState(false)
    const [countdown, setCountdown] = useState(0)

    const handleVideoFrame = useCallback((frameData: ArrayBuffer) => {
        if (jmuxerRef.current) {
            try {
                const uint8Array = new Uint8Array(frameData)
                jmuxerRef.current.feed({
                    video: uint8Array
                })
            } catch (error) {
                console.error('Error feeding video data to JMuxer:', error)
            }
        }
    }, [])

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const handleConnected = useCallback((data: any) => {
        console.log('Mirroring connected:', data)
        setIsVideoReady(true)

    }, [])

    const handleError = useCallback((error: string) => {
        console.error('Mirroring error:', error)
        setIsVideoReady(false)
    }, [])

    const {
        isConnected,
        isConnecting,
        sendTouchEvent,
        sendKeyEvent,
        connect,
        disconnect,
        error: wsError,
        screenDimensions
    } = useMirroringWebSocket({
        device,
        onConnected: handleConnected,
        onError: handleError,
        onVideoFrame: handleVideoFrame
    })

    useEffect(() => {
        if (wsError) {
            toast.error(wsError)
            handleError(wsError)
        }
    }, [wsError])

    const handleKeyCommand = useCallback((key: KeyCommand) => {
        if (isConnected && Object.values(KeyCommand).includes(key)) {
            sendKeyEvent({
                type: MessageType.KEY,
                key: key
            })
        }
    }, [isConnected, sendKeyEvent])

    const handleDisconnect = useCallback(() => {
        disconnect()
        setIsDisconnecting(true)
        setCountdown(5)

        const timer = setInterval(() => {
            setCountdown((prev) => {
                if (prev <= 1) {
                    clearInterval(timer)
                    setIsDisconnecting(false)
                    return 0
                }
                return prev - 1
            })
        }, 1000)
    }, [disconnect])

    const handleConnect = useCallback(() => {
        if (!isDisconnecting) {
            connect()
        }
    }, [connect, isDisconnecting])

    useMonitoringTouch({
        canvasRef: canvasRef as React.RefObject<HTMLCanvasElement>,
        screenWidth: screenDimensions?.width || 1080,
        screenHeight: screenDimensions?.height || 1920,
        onTouchEvent: sendTouchEvent,
        enabled: isConnected && isVideoReady
    })

    const { isFullscreen, toggleFullscreen, setFullscreenElement } = useMirroringFullscreen()

    useAspectRatio({
        canvasRef: canvasRef as React.RefObject<HTMLCanvasElement>,
        screenWidth: screenDimensions?.width,
        screenHeight: screenDimensions?.height
    })

    useEffect(() => {
        if (!videoRef.current) return

        try {
            jmuxerRef.current = new JMuxer({
                node: videoRef.current,
                mode: 'video',
                flushingTime: 10,
                fps: 60,
                debug: false,
                onReady: () => {
                    toast.success('Mirroring ready')
                },
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                onError: (error: any) => {
                    toast.error("Mirroring error:", error.message)
                }
            })

            console.log('JMuxer initialized')
        } catch (error) {
            console.error('Failed to initialize JMuxer:', error)
        }

        return () => {
            if (jmuxerRef.current) {
                try {
                    jmuxerRef.current.destroy()
                } catch (error) {
                    console.error('Error destroying JMuxer:', error)
                }
                jmuxerRef.current = null
            }
        }
    }, [isConnected])

    useEffect(() => {
        return () => {
            disconnect()
            if (jmuxerRef.current) {
                try {
                    jmuxerRef.current.destroy()
                } catch (error) {
                    console.error('Error destroying JMuxer:', error)
                    toast.error("Error destroying mirroring player")
                }
                jmuxerRef.current = null
            }
        }
    }, [disconnect, device])


    useEffect(() => {
        if (canvasRef.current && screenDimensions && mainContainerRef.current) {
            const canvas = canvasRef.current
            canvas.width = screenDimensions.width
            canvas.height = screenDimensions.height
            setFullscreenElement(mainContainerRef.current)
        }
    }, [screenDimensions, setFullscreenElement])

    if (isConnecting) {
        return (
            <div className="flex items-center justify-center bg-opacity-50 h-[300px] text-white z-20">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white mx-auto mb-2"></div>
                    <p>Connecting to {device.serial}...</p>
                </div>
            </div>
        )
    }

    if (!isConnected) {
        return (
            <div className="flex items-center justify-center bg-opacity-50 h-[300px] z-20">
                <div className="text-center">
                    <FaChromecast className='m-auto text-6xl md:text-9xl' />

                    <button
                        onClick={handleConnect}
                        disabled={isDisconnecting}
                        className={`btn ${isDisconnecting ? 'btn-disabled' : 'btn-success'}`}
                    >
                        {isDisconnecting ? (
                            <span className="flex items-center space-x-2">
                                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
                                <span>Waiting {countdown} second to complete disconnecting</span>
                            </span>
                        ) : (
                            'Connect'
                        )}
                    </button>
                    {
                        !isConnected && (
                            <div className='text-warning mt-5'>
                                <b>NOTE: </b> If you cannot connect after disconnecting or changing devices but still not connecting, please refresh the Andromodem.
                            </div>
                        )
                    }
                </div>
            </div>
        )
    }
    return (

        <>
            {
                isConnected && !isFullscreen && (
                    <div className="mb-3">
                        <MirroringTool
                            isConnected={isConnected}
                            isConnecting={isConnecting}
                            isFullscreen={isFullscreen}
                            onConnect={handleConnect}
                            onDisconnect={handleDisconnect}
                            onToggleFullscreen={toggleFullscreen}
                        />
                    </div>
                )
            }
            <div
                ref={mainContainerRef}
                className="relative flex flex-col h-[500px] md:h-[700px]"
            >
                <div className="relative flex-1 w-full min-h-0">
                    <video
                        ref={videoRef}
                        className="absolute inset-0 w-full h-full object-contain"
                        autoPlay
                        muted
                        playsInline
                    />
                    <canvas
                        ref={canvasRef}
                        className="absolute inset-0 w-full h-full object-contain z-10"
                        style={{
                            touchAction: 'none',
                            userSelect: 'none',
                            background: 'transparent',
                            pointerEvents: 'auto'
                        }}
                    />

                </div>

                {isConnected && (
                    <div className="flex-shrink-0 w-full h-[60px]">
                        <MirroringNavigation
                            isConnected={isConnected}
                            isFullscreen={isFullscreen}
                            onToggleFullscreen={toggleFullscreen}
                            onSendKeyCommand={handleKeyCommand}
                        />
                    </div>
                )}
            </div>

        </>
    )
}
export default MirroringCanvas
