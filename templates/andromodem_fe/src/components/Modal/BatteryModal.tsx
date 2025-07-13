import {FaTemperatureLow} from "react-icons/fa"
import {GiBattery100} from "react-icons/gi"
import {LuHeartPulse} from "react-icons/lu"
import {MdOutlineCategory, MdOutlineFactory} from "react-icons/md"
import {PiPlugChargingBold} from "react-icons/pi"

import type { DeviceBattery } from "@/types/device.ts"

type props = {
    battery?: DeviceBattery;
    modal_id: string;
};
const ModalBattery = ({battery, modal_id}: props) => {
    return <dialog id={modal_id} className="modal">
        <div className="modal-box">
            <h3 className="font-bold text-lg">Battery</h3>
            <ul className="list">

                <li className="list-row">
                    <div>
                        <MdOutlineFactory size={25}/>
                    </div>
                    <div>
                        <div>Charging Status</div>
                        <div className="text-xs font-semibold opacity-60">
                            {battery?.status || "-"}
                        </div>
                    </div>
                </li>
                <li className="list-row">
                    <div>
                        <GiBattery100 size={25}/>
                    </div>
                    <div>
                        <div>Level</div>
                        <div className="text-xs font-semibold opacity-60">
                            {battery?.level}
                        </div>
                    </div>
                </li>
                <li className="list-row">
                    <div>
                        <LuHeartPulse size={25}/>
                    </div>
                    <div>
                        <div>Health</div>
                        <div className="text-xs font-semibold opacity-60">
                            {battery?.health || "-"}
                        </div>
                    </div>
                </li>
                <li className="list-row">
                    <div>
                        <FaTemperatureLow size={25}/>
                    </div>
                    <div>
                        <div>Battery Temperature</div>
                        <div className="text-xs font-semibold opacity-60">
                            {battery?.temperature} °C
                        </div>
                    </div>
                </li>
                <li className="list-row">
                    <div>
                        <PiPlugChargingBold size={25}/>
                    </div>
                    <div>
                        <div>Charging Counter</div>
                        <div className="text-xs font-semibold opacity-60">
                            {battery?.charge_counter || "-"}
                        </div>
                    </div>
                </li>
                <li className="list-row">
                    <div>
                        <MdOutlineCategory size={25}/>
                    </div>
                    <div>
                        <div>Technology</div>
                        <div className="text-xs font-semibold opacity-60">
                            {battery?.technology || "-"}
                        </div>
                    </div>
                </li>
            </ul>
            <div className="modal-action">
                <form method="dialog">
                    <button className="btn btn-outline">Close</button>
                </form>
            </div>
        </div>
    </dialog>
}
export default ModalBattery