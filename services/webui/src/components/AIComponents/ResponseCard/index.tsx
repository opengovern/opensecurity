import { FunctionComponent, useEffect, useState } from 'react';
import KTable from '../Table';
import { useAnimatedText } from '../../../utilities/useAnimate';
import LoadingDots from '../Loading';

interface ChatCardProps {
    response: any
    key: number
    loading: boolean
    error: string
    scroll: Function
    ref: any
    time: number
    text: string
    suggestions: string[] | undefined
    onClickSuggestion: Function
    id: string
    isWelcome?: boolean
    date: string
    clarify_needed: boolean
    clarify_questions: ClarifyQuestion[]
    chat_id: string
    pre_loaded: boolean
}
interface ClarifyQuestion {
  question : string

}

const KResponseCard: FunctionComponent<any> = ({
  response,
  key,
  loading,
  error,
  scroll,
  ref,
  time,
  text,
  suggestions,
  onClickSuggestion,
  pre_loaded,
  isWelcome,
  id,
  date,
  clarify_needed,
  clarify_questions,
  chat_id,
}: ChatCardProps) => {
  useEffect(() => {
    scroll();
    if (id && id != '' && clarify_needed == false) {
      setText('Running sql query ');
    } else {
      setText('Interpreting');
    }
  }, [error, loading, response, id]);

  const [loadingText, setText] = useState('Running sql query');
  const [showText, setShowText] = useState(' ');
  const animatedtText = useAnimatedText(showText, 1);
  useEffect(() => {
    setShowText(text);
  }, [text]);

  return (
      <>
          <div
              key={key}
              className="    flex-col gap-2   flex justify-start items-start max-w-[95%]"
          >
              <div className="flex flex-row gap-1 justify-start items-center">
                  <span>👋</span>
                  <span className="text-base">Assistant</span>
                  <span className=" text-sm text-slate-800 dark:text-slate-300">
                      {date}
                  </span>
              </div>
              {loading || (!response && !text) ? (
                  <>
                      <div className="flex flex-row gap-2 justify-start items-center w-full">
                          <span className=" text-slate-800 dark:text-slate-200 w-fit max-w-[95%]">
                              {pre_loaded
                                  ? loadingText
                                  : useAnimatedText(loadingText, 1).text}
                          </span>
                          <LoadingDots />
                      </div>
                  </>
              ) : (
                  <>
                      {response ? (
                          <>
                              {' '}
                              <div className="flex flex-col gap-4 h-full  w-full justify-start items-start">
                                  {error && error != '' ? (
                                      <>
                                          <div className="rounded-3xl    p-2 px-4 my-2 text-red-700 text-left sm:w-fit w-full flex flex-row gap-2">
                                              {/* <RiAlertLine /> */}
                                              {pre_loaded
                                                  ? error
                                                  : useAnimatedText(error, 3)
                                                        .text}
                                          </div>
                                      </>
                                  ) : (
                                      <>
                                          <div className="flex flex-col gap-4 h-full  w-full justify-start items-start">
                                              {clarify_needed ? (
                                                  <>
                                                      <span className="text-slate-800 dark:text-slate-200">
                                                          Please answer below
                                                          questions to clarify
                                                          query.{' '}
                                                      </span>

                                                      <span className="text-slate-800 dark:text-slate-200">
                                                          {pre_loaded
                                                              ? clarify_questions?.map(
                                                                    (
                                                                        question
                                                                    ) => {
                                                                        return question?.question
                                                                    }
                                                                ).join (",")
                                                              : useAnimatedText(
                                                                    clarify_questions
                                                                        ?.map(
                                                                            (
                                                                                question: any
                                                                            ) => {
                                                                                return question?.question
                                                                            }
                                                                        )
                                                                        .join(
                                                                            ','
                                                                        ),
                                                                    2
                                                                ).text}
                                                      </span>
                                                  </>
                                              ) : (
                                                  <>
                                                      <span className="text-slate-800 dark:text-slate-200">
                                                          {pre_loaded
                                                              ? showText
                                                              : useAnimatedText(
                                                                    showText,
                                                                    3
                                                                ).text}
                                                      </span>
                                                  </>
                                              )}
                                              {!pre_loaded &&
                                                  animatedtText.done && (
                                                      <>
                                                          <KTable
                                                              ref={ref}
                                                              scroll={scroll}
                                                              result={response}
                                                              key={
                                                                  key + 'table'
                                                              }
                                                              time={time}
                                                              suggestions={
                                                                  suggestions
                                                              }
                                                              onClickSuggestion={
                                                                  onClickSuggestion
                                                              }
                                                              isWelcome={
                                                                  isWelcome
                                                              }
                                                              chat_id={chat_id}
                                                          />
                                                      </>
                                                  )}
                                              {pre_loaded && (
                                                  <>
                                                      <KTable
                                                          ref={ref}
                                                          scroll={scroll}
                                                          pre_loaded={
                                                              pre_loaded
                                                          }
                                                          result={response}
                                                          key={key + 'table'}
                                                          time={time}
                                                          suggestions={
                                                              suggestions
                                                          }
                                                          onClickSuggestion={
                                                              onClickSuggestion
                                                          }
                                                          isWelcome={isWelcome}
                                                          chat_id={chat_id}
                                                      />
                                                  </>
                                              )}
                                          </div>
                                      </>
                                  )}
                              </div>
                          </>
                      ) : (
                          <>
                              <div className="rounded-3xl bg-gray-800   p-2 px-4 my-2 text-red-500 text-center w-fit flex flex-row gap-2">
                                  {/* <RiAlertLine /> */}
                                  Can not get a results
                              </div>
                          </>
                      )}
                  </>
              )}
          </div>
      </>
  )
};

export default KResponseCard;
