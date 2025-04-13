import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import { ChangeEvent, ChangeEventHandler, FunctionComponent, KeyboardEvent, useContext, useEffect, useRef, useState, } from "react";

interface InputProps {
  onChange: ChangeEventHandler;
  value: string;
  onSend: Function;
  disabled: boolean;
  chats: any;
}

const KInput: FunctionComponent<any> = ({
  onChange,
  value,
  onSend,
  disabled,
  chats
}: InputProps) => {
  // const { Theme, changeTheme } = useContext(ThemeContext);
  const [is_dark, setIsDark] = useState(false);


  const handleKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter" && event.shiftKey) {
      // Handle Shift+Enter: Add a new line
      event.preventDefault(); // Prevent the default behavior of form submission
      const e = {
        target: {
          value: value + "\n",
        },
      };
      onChange(e as ChangeEvent<HTMLInputElement>); // Update the input value with a new line
    } else if (event.key === "Enter") {
      // Handle Enter: Send the input
      event.preventDefault(); // Prevent default behavior
      handleSend();
    }
  };

  const handleSend = () => {
    if (value.trim() !== "") {
      const trimmedValue = value.trim(); // Remove leading and trailing whitespaces
      const e = {
        target: {
          value: trimmedValue,
        },
      };
      onChange(e as ChangeEvent<HTMLInputElement>); // Update the input value with the trimmed value
      onSend(); // Trigger the send callback with the current input value
    }
   
  };
  const firstRender= useRef(false);
  // useEffect(()=>{
  //   if(!firstRender.current || chats.length==0){
  //      if (value == '' || value == undefined) {
  //        // @ts-ignore
  //        onChange({ target: { value: 'Get List of all identities' } });
  //      }
  //       firstRender.current=true;

  //   }
   
  // },[value])
  return (
    <>
      <div
        className={` #dark:bg-gray-950 #bg-slate-200 w-full px-1 sm:px-0 flex justify-center items-center flex-row  max-w-[90%] absolute bottom-0    ${window.innerWidth<640 && 'left-0'}  `}
      >
        <div className="   sm:w-full w-[98%] h-12 flex flex-row sm:justify-start justify-center items-center bg-gray-300 dark:bg-gray-950  rounded-full  dark:border-slate-600 border-slate-200 border ">
          <input
            value={value}
            onKeyDown={(e) => {
              if (!disabled) {
                handleKeyDown(e);
              }
            }}
            onChange={onChange}
            placeholder="Ask from AI"
            className="w-full  max-w-[94%] k-input bg-gray-300 dark:bg-gray-950 border-none h-10  rounded-full text-slate-900 dark:text-slate-300 shadow-none placeholder:text-slate-600 dark:placeholder:text-slate-500 focus:shadow-none focus:outline-0 focus:ring-0     resize-none overflow-y-auto text-base p-2 px-4  "
          ></input>
          <button
            onClick={() => {
              onSend();
            }}
            disabled={disabled}
            className={` absolute  right-3 ${is_dark ? 'text-white' : 'text-black'} rounded-full p-2 disabled:text-gray-500 `}
          >
            <PaperAirplaneIcon color="black" className="w-5" />
          </button>
        </div>
      </div>
    </>
  );
};

export default KInput;
