import React from 'react'
import { TbArrowsMaximize } from 'react-icons/tb'

interface MirroringToolProps {
  isConnected: boolean;
  isConnecting: boolean;
  isFullscreen: boolean;
  onConnect: () => void;
  onDisconnect: () => void;
  onToggleFullscreen: () => void;
}

const MirroringTool: React.FC<MirroringToolProps> = ({
  onDisconnect,
  onToggleFullscreen
}) => {
  return (
    <>
        <button
          onClick={onDisconnect}
          className="btn btn-soft btn-sm md:btn-md btn-error"
        >
          Disconnect
        </button>

      <button
        onClick={onToggleFullscreen}
        className="btn btn-primary ml-3 btn-sm md:btn-md"
        title="Enter Fullscreen"
      >
        <TbArrowsMaximize fontSize={24} />
      </button>
      </>
      
  )
}
export default MirroringTool