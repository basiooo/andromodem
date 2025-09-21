import React from 'react'

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
          className="btn btn-soft btn-error"
        >
          Disconnect
        </button>

      <button
        onClick={onToggleFullscreen}
        className="btn btn-primary"
        title="Enter Fullscreen"
      >
        ⤢
      </button>
      </>
      
  )
}
export default MirroringTool