import { FunctionComponent } from "react";

interface ChatCardProps {
  message: string;
  keyNumber: number;
  date : string
}

const KChatCard: FunctionComponent<any> = ({
  message,
  keyNumber,
  date
}: ChatCardProps) => {
  return (
    <div key={keyNumber} className="  w-full flex flex-col   justify-start items-start">
      <div className="flex flex-row gap-1 justify-start items-center">
        <span>You</span>
        <span className=" text-sm text-slate-800 dark:text-slate-300">{date}</span>
      </div>
      <span className=" rounded-3xl    text-slate-800 dark:text-slate-200 text-left w-fit">{message}</span>
    </div>
  );
};

export default KChatCard;
