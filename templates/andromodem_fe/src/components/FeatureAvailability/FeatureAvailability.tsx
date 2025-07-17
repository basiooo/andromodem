import { type FC } from "react" 

import useFeatureAvailability from "@/hooks/useFeatureAvailability"
import { type Device,DeviceState } from "@/types/device"
import { showModal } from "@/utils/common"

import FeatureAvailabilityModal from "@/components/Modal/FeatureAvailabilityModal"

const FeatureAvailability: FC<{ device: Device }> = ({ device }) => {
    const { deviceFeatureAvailabilities } = useFeatureAvailability(device)
    if (deviceFeatureAvailabilities === null) {
        return null
    }
    return (
        <div>
            <FeatureAvailabilityModal features={deviceFeatureAvailabilities} modal_id="feature_availability_modal" />
            <button disabled={device.state !== DeviceState.ONLINE} className="btn btn-xs btn-info" onClick={
                () => showModal("feature_availability_modal")
            }>Feature Availability</button>
        </div>
    )
}

export default FeatureAvailability
