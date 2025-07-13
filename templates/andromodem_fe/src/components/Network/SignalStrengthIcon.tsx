import type { FC } from "react"
import {
    MdSignalCellular1Bar,
    MdSignalCellular2Bar,
    MdSignalCellular3Bar,
    MdSignalCellular4Bar
} from "react-icons/md"

const SignalStrengthIcon: FC<{ level: number }> = ({ level }) => {
    console.log(level)
    switch (level) {
        case 1:
            return <MdSignalCellular1Bar fontSize={30} />
        case 2:
            return <MdSignalCellular2Bar fontSize={30} />
        case 3:
            return <MdSignalCellular3Bar fontSize={30} />
        case 4:
            return <MdSignalCellular4Bar fontSize={30} />
        default:
            return "Unknown"
    }
}

export default SignalStrengthIcon