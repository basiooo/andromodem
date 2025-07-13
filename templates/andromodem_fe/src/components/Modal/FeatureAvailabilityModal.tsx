import { FaCheck } from "react-icons/fa6"
import { MdBlock } from "react-icons/md"

import type { FeatureAvailabilities } from "@/types/device.ts"

type props = {
    features: FeatureAvailabilities;
    modal_id: string;
};
const FeatureAvailabilityModal = ({features, modal_id}: props) => {
    return <dialog id={modal_id} className="modal">
        <div className="modal-box">
            <h3 className="font-bold text-lg">Feature Availability</h3>
            <p className="text-sm text-gray-500 mb-4">
    This is a list of features reported as available on the selected device. Actual availability may vary depending on device conditions or system limitations.
</p>
            <table className="table table-zebra text-xs">
                <thead>
                    <tr>
                        <th>Feature</th>
                        <th>Status</th>
                        <th>Message</th>
                    </tr>
                </thead>
                <tbody>
                    {features.map((feature) => (
                        <tr key={feature.key}>
                            <td>{feature.feature}</td>
                            <td>{feature.available ? <FaCheck color="green" /> : <MdBlock color="red" />}</td>
                            <td>{feature.message}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <div className="modal-action">
                <form method="dialog">
                    <button className="btn btn-outline">Close</button>
                </form>
            </div>
        </div>
    </dialog>
}
export default FeatureAvailabilityModal