import type { FC } from "react"
import {
    ScreenResolution,
    type ScreenResolutionValue,
    BitRate,
    type BitRateValue,
    FPS,
    type FPSValue
} from "@/types/mirroring"

type MirroringConfigProps = {
    screenResolution: ScreenResolutionValue
    setScreenResolution: (value: ScreenResolutionValue) => void
    bitrate: BitRateValue
    setBitrate: (value: BitRateValue) => void
    fps: FPSValue
    setFps: (value: FPSValue) => void
    handleConnect: () => void
    isDisconnecting: boolean
    countdown: number
    isConnected: boolean
}

const MirroringConfig: FC<MirroringConfigProps> = ({
    screenResolution,
    setScreenResolution,
    bitrate,
    setBitrate,
    fps,
    setFps,
    handleConnect,
    isDisconnecting,
    countdown,
    isConnected
}) => {
    return (
        <div className="text-center">
            <div className="mb-5">
                <b className="text-2xl">Mirroring Configuration</b>
            </div>
            <div className="form-control text-left max-w-xs m-auto">
                <label className="label">
                    <span className="label-text">Screen Resolution</span>
                </label>
                <select
                    value={screenResolution}
                    disabled={countdown > 0}
                    onChange={(e) => setScreenResolution(Number(e.target.value) as ScreenResolutionValue)}
                    className="select select-sm select-bordered"
                >
                    {Object.values(ScreenResolution).map((resolution) => (
                        <option key={resolution} value={resolution}>
                            {resolution} Pixel
                        </option>
                    ))}
                </select>
            </div>

            <div className="form-control text-left max-w-xs m-auto">
                <label className="label">
                    <span className="label-text">BitRate</span>
                </label>
                <select
                    value={bitrate}
                    disabled={countdown > 0}
                    onChange={(e) => setBitrate(Number(e.target.value) as BitRateValue)}
                    className="select select-sm select-bordered"
                >
                    {Object.entries(BitRate).map(([key, value]) => (
                        <option key={value} value={value}>
                            {key} Mbps
                        </option>
                    ))}
                </select>
            </div>

            <div className="form-control text-left max-w-xs m-auto">
                <label className="label">
                    <span className="label-text">FPS</span>
                </label>
                <select
                    value={fps}
                    disabled={countdown > 0}
                    onChange={(e) => setFps(Number(e.target.value) as FPSValue)}
                    className="select select-sm select-bordered"
                >
                    {Object.values(FPS).map((value) => (
                        <option key={value} value={value}>
                            {value} FPS
                        </option>
                    ))}
                </select>
            </div>

            <button
                onClick={handleConnect}
                disabled={isDisconnecting}
                className={`btn mt-5 ${isDisconnecting ? "btn-disabled" : "btn-success"}`}
            >
                {isDisconnecting ? (
                    <span className="flex items-center space-x-2">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
                        <span>Waiting {countdown} second to complete disconnecting</span>
                    </span>
                ) : (
                    "Connect"
                )}
            </button>

            {!isConnected && (
                <div className="text-warning m-3">
                    <b>NOTE: </b> If you cannot connect after disconnecting or changing devices but still not connecting, please refresh the Andromodem.
                </div>
            )}
        </div>
    )
}

export default MirroringConfig