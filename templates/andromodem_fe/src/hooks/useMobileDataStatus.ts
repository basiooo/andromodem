import { MobileDataConnectionState, type MobileDataConnectionStateValue, type SimInfo } from "@/types/network"

const useMobileDataStatus = () => {
  const isMobileDataEnabled = (connectionState: MobileDataConnectionStateValue) => {
    return (
      connectionState === MobileDataConnectionState.CONNECTED || 
      connectionState === MobileDataConnectionState.CONNECTING
    )
  }

  const hasMobileDataEnabled = (sims: SimInfo[]) => {
    return sims?.some((sim) => isMobileDataEnabled(sim.connection_state))
  }

  return {
    isMobileDataEnabled,
    hasMobileDataEnabled
  }
}

export default useMobileDataStatus