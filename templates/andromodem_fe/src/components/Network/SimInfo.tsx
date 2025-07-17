import type { FC } from "react"
import { TbMobiledata } from "react-icons/tb"

import useMobileDataStatus from "@/hooks/useMobileDataStatus"
import type { SimInfo as SimInfoType } from "@/types/network"
import SignalStrengthIcon from "@/components/Network/SignalStrengthIcon"


interface SimInfoProps {
  sims: SimInfoType[]
  activeTab: string
  onTabChange: (tabKey: string) => void
}

const SimInfoDisplay: FC<SimInfoProps> = ({ sims, activeTab, onTabChange }) => {
  const { isMobileDataEnabled } = useMobileDataStatus()

  if (!sims?.length) {
    return (
      <div className="card w-full bg-base-100 shadow-sm">
        <div className="card-body">
          <h1 className="text-2xl text-center font-bold my-5 text-red-500">Cannot extract sim info</h1>
        </div>
      </div>
    )
  }

  return (
    <div className="card w-full bg-base-100 shadow-sm">
      <div className="card-body">
        <div role="tablist" className="tabs tabs-lifted">
          {
            sims.map((sim) => (
              <>
                <button 
                  type="button" 
                  role="tab" 
                  className={`tab flex items-center ${activeTab === sim.sim_slot.toString() ? "tab-active !bg-base-200 text-lg" : ""}`}
                  onClick={() => onTabChange(sim.sim_slot.toString())}
                >
                  {sim.name}
                  {
                    isMobileDataEnabled(sim.connection_state) ?
                      <TbMobiledata color="green" fontSize={23} /> :
                      ""
                  }
                </button>
                <div role="tabpanel" className="tab-content bg-base-200 border-base-300 rounded-box p-6 overflow-auto ">
                  <table className="table">
                    <tbody>
                      <tr>
                        <td className='text-xs md:text-sm'>Operator Name</td>
                        <td className='text-xs md:text-sm'>{sim.name}</td>
                      </tr>
                      <tr>
                        <td className='text-xs md:text-sm'>Network Type</td>
                        <td className='text-xs md:text-sm'>{Object.keys(sim.signal_strength).length === 0 ? "Unknown" : Object.keys(sim.signal_strength)[0]}  </td>
                      </tr>
                      <tr>
                        <td className='text-xs md:text-sm'>Sim Slot</td>
                        <td className='text-xs md:text-sm'>{sim.sim_slot}</td>
                      </tr>
                      <tr>
                        <td className='text-xs md:text-sm'>Mobile Data Status</td>
                        <td className='text-xs md:text-sm'>{sim.connection_state.length > 0 ? sim.connection_state : "Unknown"}</td>
                      </tr>
                      <tr>
                        <td className='text-xs md:text-sm'>Signal Strength</td>
                        <td className='text-xs md:text-sm'>
                          <SignalStrengthIcon level={Number(sim.signal_strength[Object.keys(sim.signal_strength)[0]]?.level)} />
                        </td>
                      </tr>
                      {
                        Object.keys(sim.signal_strength).length > 0 && Object.keys(sim.signal_strength[Object.keys(sim.signal_strength)[0] as keyof typeof sim.signal_strength]).length > 0 ?
                          Object.entries(Object.values(sim.signal_strength)[0] as Record<string, string>).map(([k, v]) => (
                            <tr key={k}>
                              <td className='text-xs md:text-sm '>{k}</td>
                              <td className='text-xs md:text-sm '>{v}</td>
                            </tr>
                          )) :
                          <>
                            <tr>
                              <td colSpan={2} className='text-center text-red-600'>
                                <div>
                                  <p>Unable to extract signal strength</p>
                                  <p className='text-center'>This feature is only available on Android 10 or above.</p>
                                </div>
                              </td>
                            </tr>
                          </>
                      }
                    </tbody>
                  </table>
                </div>
              </>
            ))
          }
        </div>
      </div>
    </div>
  )
}

export default SimInfoDisplay