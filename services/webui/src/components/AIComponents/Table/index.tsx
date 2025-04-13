import { FunctionComponent, useEffect, useState } from "react";
import Table from "@cloudscape-design/components/table";
import {
  AppLayout,
  Box,
  ExpandableSection,
  Header,
  Modal,
  Pagination,
  Popover,
  SpaceBetween,
  SplitPanel,
  Tabs,
} from "@cloudscape-design/components";
import axios from "axios";
import Tooltip from "../Tooltip";
import CustomPagination from "../../Pagination";
import { useAnimatedText } from "../../../utilities/useAnimate";
import { EpochtoSecond } from "../../../utilities/dateDisplay";

export const capitalizeFirstLetters = (string: string) => {
  
  const splitStr = string.toLowerCase().split(" ");
  for (let i = 0; i < splitStr.length; i += 1) {
    // You do not need to check if i is larger than splitStr length, as your for does that for you
    // Assign it back to the array
    splitStr[i] =
      splitStr[i].charAt(0).toUpperCase() + splitStr[i].substring(1);
  }
  // Directly return the joined string
  return splitStr.join(" ");
};


export const snakeCaseToLabel = (string: string) =>
  capitalizeFirstLetters(
    string
      .toLowerCase()
      .replace(/([-_][a-z])/g, (group) => group.replace("_", " "))
  );


export const getTable = (
  headers: string[] | undefined,
  details: any[][] | undefined
) => {
  const columns: any[] = [];
  const rows: any[] = [];
  const headerField = headers?.map((value, idx) => {
    if (headers.filter((v) => v === value).length > 1) {
      return `${value}-${idx}`;
    }
    return value;
  });
  if (headers && headers.length) {
    for (let i = 0; i < headers.length; i += 1) {
      const isHide = headers[i][0] === "_";
  
      columns.push(snakeCaseToLabel(headers[i]));
  
    }
  }
  if (details && details.length) {
    for (let i = 0; i < details.length; i += 1) {
      const row: any = {};
      for (let j = 0; j < columns.length; j += 1) {
        row[headerField?.at(j) || ''] =
          typeof details[i][j] == 'string'
            ? // @ts-ignore
              details[i][j]
            : // @ts-ignore
              typeof details[i][j] == 'number' && details[i][j] % 1 !== 0
              ? details[i][j].toFixed(2)
              : JSON.stringify(details[i][j]);
       
      }
      rows.push(row);
    }
  }
  const count = rows.length;

  return {
    columns,
    rows,
    count,
  };
};
export const getTableCloudScape = (
  headers: string[] | undefined,
  details: any[][] | undefined
) => {
  const columns: any[] = [];
  const rows: any[] = [];
  const column_def: any[] = [];
  const headerField = headers?.map((value, idx) => {
    if (headers.filter((v) => v === value).length > 1) {
      return `${value}-${idx}`;
    }
    return value;
  });
  if (headers && headers.length) {
    for (let i = 0; i < headers.length; i += 1) {
      const isHide = headers[i][0] === "_";
  
      columns.push({
        id: headerField?.at(i),
        header: (
          <>
            <div className="pl-2">{snakeCaseToLabel(headers[i])}</div>
          </>
        ),
        // @ts-ignore
        cell: (item: any) => (
          <div className="pl-2">
            {/* @ts-ignore */}
            {typeof item[headerField?.at(i)] == "string"
              ? // @ts-ignore
                item[headerField?.at(i)]
              : // @ts-ignore
                JSON.stringify(item[headerField?.at(i)])}
          </div>
        ),
        maxWidth: "200px",
        // sortingField: 'id',
        // isRowHeader: true,
        // maxWidth: 150,
      });
      column_def.push({
        id: headerField?.at(i),
        visible: !isHide,
      });
    }
  }
  if (details && details.length) {
    for (let i = 0; i < details.length; i += 1) {
      const row: any = {};
      for (let j = 0; j < columns.length; j += 1) {
        row[headerField?.at(j) || ""] = details[i][j];
        //     typeof details[i][j] === 'string'
        //         ? details[i][j]
        //         : JSON.stringify(details[i][j])
      }
      rows.push(row);
    }
  }
  const count = rows.length;

  return {
    columns,
    column_def,
    rows,
    count,
  };
};


interface ChatCardProps {
  result: any;
  key: number;
  scroll: Function;
  ref: any;
  time: number;
  suggestions: string[] | undefined;
  onClickSuggestion: Function;
  isWelcome?: boolean;
  chat_id: string;
  pre_loaded: boolean
}

const KTable: FunctionComponent<any> = ({
  result,
  key,
  scroll,
  ref,
  time,
  suggestions,
  onClickSuggestion,
  isWelcome,
  chat_id,
  pre_loaded,
}: ChatCardProps) => {
  const [page, setPage] = useState(0);
  const [open, setOpen] = useState(false);

const Downloadchats = () => {
 
   let url = ''
   if (window.location.origin === 'http://localhost:3000') {
       url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
   } else {
       url = window.location.origin
   }
   // @ts-ignore
   const token = JSON.parse(localStorage.getItem('openg_auth')).token

   const config = {
       headers: {
           Authorization: `Bearer ${token}`,
       },
       responseType: "blob", // Ensure we get a file
   }
  

  axios
//   @ts-ignore
    .get(`${url}/main/core/api/v4/chatbot/chat/${chat_id}/download` , config)
    .then((res) => {
      // @ts-ignore
       let fileName = `export_${chat_id}.csv`; // Default filename

       

       // Create a download link
       const url = window.URL.createObjectURL(new Blob([res.data]));
       const link = document.createElement('a');
       link.href = url;
       link.setAttribute('download', fileName); // Set filename dynamically
       document.body.appendChild(link);
       link.click();
       document.body.removeChild(link);
    })
    .catch((err: any) => {
      scroll();
    });
};

  useEffect(() => {
    scroll();
    console.log(result)
  }, [result]);
  return (
      <div
          className={`flex flex-col ${
              !isWelcome && 'gap-4'
          } h-full w-full  justify-between items-start`}
      >
          {!isWelcome && (
              <>
                  {getTable(result.headers, result?.result).count !== 0 ? (
                      <>
                          <div
                              key={key}
                              ref={ref}
                              className="   flex justify-start items-start max-h-[50vh]  w-full"
                          >
                              <span className="       my-2 text-slate-200 text-center w-full    overflow-x-auto">
                                  <table className="table-auto w-full border-slate-500 p-4  rounded-t-2xl bg-gray-700    border-collapse   ">
                                      <thead className="mb-2 rounded-xl w-full ">
                                          <tr className="     ">
                                              {/* @ts-ignore */}

                                              {getTable(
                                                  result.headers,
                                                  result?.result
                                              ).columns?.map((col: any) => {
                                                  return (
                                                      <>
                                                          <th className="text-white text-left truncate      p-2 sm:p-4">
                                                              {col}
                                                          </th>
                                                      </>
                                                  )
                                              })}
                                          </tr>
                                      </thead>
                                      <tbody className="w-full">
                                          {/* @ts-ignore */}
                                          {getTable(
                                              result.headers,
                                              result?.result
                                          )
                                              .rows.slice(0 * 5, (0 + 1) * 5)
                                              ?.map((item: any, index: any) => {
                                                  return (
                                                      <tr
                                                          className={` text-left  ${
                                                              index <
                                                              // @ts-ignore
                                                              getTable(
                                                                  result.headers,
                                                                  result?.result
                                                              )?.rows.slice(
                                                                  0 * 5,
                                                                  (0 + 1) * 5
                                                              ).length -
                                                                  1
                                                                  ? ' border-b border-slate-400 '
                                                                  : ''
                                                          }  bg-gray-950`}
                                                      >
                                                          {Object.keys(
                                                              item
                                                          ).map((key) => {
                                                              return (
                                                                  <td className="text-white max-w-28 truncate    p-2 sm:p-4">
                                                                      {
                                                                          item[
                                                                              key
                                                                          ]
                                                                      }
                                                                  </td>
                                                              )
                                                          })}
                                                      </tr>
                                                  )
                                              })}
                                      </tbody>
                                  </table>
                                  <div className="w-full flex flex-row justify-end gap-4 items-center  bg-gray-700 p-2 rounded-b-2xl ">
                                      {result?.result &&
                                          getTable(
                                              result.headers,
                                              result?.result
                                          )?.count > 5 && (
                                              <>
                                                  <div
                                                      className="flex flex-row gap-2  text-slate-200 hover:bg-slate-500 px-2   justify-end items-end w-fit cursor-pointer"
                                                      onClick={(e) => {
                                                          setOpen(true)
                                                      }}
                                                  >
                                                      <span className=" text-center min-w-max">
                                                          See all{' '}
                                                          {
                                                              getTable(
                                                                  result.headers,
                                                                  result?.result
                                                              )?.count
                                                          }{' '}
                                                          results
                                                      </span>
                                                      {/* <RiDropdownList /> */}
                                                  </div>
                                                  <div
                                                      className="flex flex-row gap-2  text-slate-200 rounded-xl hover:bg-slate-500 px-2    justify-end items-end w-fit cursor-pointer"
                                                      onClick={(e) => {
                                                          Downloadchats()
                                                      }}
                                                  >
                                                      <span className=" text-center min-w-max">
                                                          Download{' '}
                                                          {
                                                              getTable(
                                                                  result.headers,
                                                                  result?.result
                                                              )?.count
                                                          }{' '}
                                                          results
                                                      </span>
                                                      {/* <RiDownloadLine /> */}
                                                  </div>
                                              </>
                                          )}
                                  </div>
                              </span>
                          </div>
                      </>
                  ) : (
                      <>
                          <div className="rounded-3xl dark:bg-gray-800   p-2 px-4 my-2  dark:text-yellow-400 text-yellow-800 text-center w-fit flex flex-row gap-2">
                              {/* <RiErrorWarningLine /> */}
                              No results.
                          </div>
                      </>
                  )}
              </>
          )}
          <div className="  sm:grid sm:grid-cols-12 sm:max-w-[98%]  flex flex-col gap-2 justify-start items-center w-full    ">
              {suggestions &&
                  suggestions.length > 0 &&
                  suggestions?.map((suggestion: string) => {
                      return (
                          <>
                              <Tooltip
                                  text={suggestion}
                                  className=" col-span-4"
                              >
                                  <div
                                      onClick={() => {
                                          onClickSuggestion(suggestion)
                                      }}
                                      className=" rounded-3xl bg-slate-400 hover:bg-slate-600 flex flex-row gap-2  cursor-pointer   p-2 px-4 my-2 text-slate-800 text-center "
                                  >
                                      <span
                                          className={`truncate ${
                                              isWelcome ? 'w-full' : 'w-full'
                                          } `}
                                      >
                                          {' '}
                                          {pre_loaded
                                              ? suggestion
                                              : useAnimatedText(suggestion, 4)
                                                    .text}
                                      </span>
                                      {/* <RiSparklingLine /> */}
                                  </div>
                              </Tooltip>
                          </>
                      )
                  })}
          </div>
          {result?.result && !isWelcome && (
              <>
                  <div
                      className="flex flex-row gap-2 text-sm dark:text-slate-200 w-fit  justify-end items-center "
                      onClick={(e) => {
                          // setOpen(true);
                      }}
                  >
                      <span>Took {EpochtoSecond(time)}s</span>
                      {/* <RiTable2 /> */}
                  </div>
              </>
          )}
          <Modal
              visible={open}
              onDismiss={() => {
                  setOpen(false)
              }}
          >
              <div>
                  <div
                      key={key}
                      className="   flex justify-start items-start flex-col max-w-[40dvw] my-2 "
                  >
                      <Table
                          className="     p-4 dark:bg-gray-700 custom-table   "
                          // resizableColumns
                          variant="full-page"
                          renderAriaLive={({
                              firstIndex,
                              lastIndex,
                              totalItemsCount,
                          }) =>
                              `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                          }
                          onSortingChange={(event) => {
                              // setSort(event.detail.sortingColumn.sortingField)
                              // setSortOrder(!sortOrder)
                          }}
                          // sortingColumn={sort}
                          // sortingDescending={sortOrder}
                          // sortingDescending={sortOrder == 'desc' ? true : false}
                          // @ts-ignore
                          // stickyHeader={true}
                          resizableColumns={true}
                          // stickyColumns={
                          //  {   first:1,
                          //     last: 1}
                          // }
                          onRowClick={(event) => {}}
                          columnDefinitions={
                              getTableCloudScape(
                                  result?.headers,
                                  result?.result
                              ).columns
                          }
                          columnDisplay={
                              getTableCloudScape(
                                  result?.headers,
                                  result?.result
                              ).column_def
                          }
                          enableKeyboardNavigation
                          // @ts-ignore
                          items={getTableCloudScape(
                              result?.headers,
                              result?.result
                          ).rows?.slice(page * 10, (page + 1) * 10)}
                          loading={false}
                          loadingText="Loading resources"
                          // stickyColumns={{ first: 0, last: 1 }}
                          // stripedRows
                          trackBy="id"
                          empty={
                              <Box
                                  margin={{
                                      vertical: 'xs',
                                  }}
                                  textAlign="center"
                                  color="inherit"
                              >
                                  <SpaceBetween size="m">
                                      <b>No Results</b>
                                  </SpaceBetween>
                              </Box>
                          }
                          header={
                              <Header
                                  className="w-full"
                                  actions={
                                      <CustomPagination
                                          currentPageIndex={page + 1}
                                          onChange={({ detail }: any) => {
                                              setPage(
                                                  detail.currentPageIndex - 1
                                              )
                                          }}
                                          pagesCount={Math.ceil(
                                              getTableCloudScape(
                                                  result.headers,
                                                  result?.result
                                              ).count / 10
                                          )}
                                      />
                                  }
                              >
                                  Results{' '}
                                  <span className=" font-medium">
                                      (
                                      {
                                          getTableCloudScape(
                                              result?.headers,
                                              result?.result
                                          ).count
                                      }
                                      )
                                  </span>
                              </Header>
                          }
                      />
                  </div>
              </div>
          </Modal>
      </div>
  )
};

export default KTable;
