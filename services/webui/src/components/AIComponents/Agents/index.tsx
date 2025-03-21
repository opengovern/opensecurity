import axios from 'axios';
import { useEffect, useState } from 'react';
import { Agent } from './types';
import { useNavigate } from 'react-router-dom';
import LoadingDots from '../Loading';
import Tooltip from '../Tooltip';
import { Button, Modal } from '@cloudscape-design/components'
import Cal, { getCalApi } from '@calcom/embed-react'
import { Flex } from '@tremor/react'
function Agents() {
      const [open,setOpen] = useState(false)
      const [openCal, setOpenCal] = useState(false)
  
  const [agents, setAgents] = useState<Agent[]>([
      {
          name: 'Identity & Access',
          description:
              'Delivers information on user identities, access controls, and activity.',
          welcome_message:
              'Hi there! This is your Identity & Access Agent. I can help you with anything related to identity management and access tools. What can I assist you with today? For example, you can ask me things like:',
          sample_questions: [
              'Get me the list of users who have access to Azure Subscriptions.',
              'Get me all SPNs with expired passwords.',
              'Show me the access activity for user John Doe.',
          ],
          id: 'identity_access',
          enabled: true,
          is_available: true,
      },
      {
          name: 'DevOps',
          description:
              'Provides data on secure code, deployment, and automation workflows.',
          welcome_message:
              'Hello! This is your DevOps Agent. I can provide insights into secure code, deployment, and automation workflows. How can I assist you today? For instance, you could ask me:',
          sample_questions: [
              'What are the latest secure code scan results?',
              'Show me the deployment status for the production environment.',
              'Provide a report on automated workflow execution times.',
          ],
          id: 'devops',
          enabled: true,
          is_available: true,
      },
      {
          name: 'Sales',
          description:
              'Answers questions about sales activities, including: total activities per rep, activity breakdown by type, activity per deal, deal stage progression activity, and time to first activity. Enables analysis of sales rep performance and deal progression.',
          welcome_message:
              'Hello! This is your Sales Agent. I can answer questions about sales activities, including total activities per rep, activity breakdown by type, activity per deal, deal stage progression activity, and time to first activity. What can I help you with? For example:',
          sample_questions: [
              'Show me the total activities (calls, emails, meetings) for each sales rep.',
              'Which sales reps had the most/least activities?',
              'How much activity was logged for each deal we closed?',
          ],
          id: 'sales',
          enabled: true,
          is_available: true,
      },
  ])
  const selected_agent = {
      id: 'identity_access',
  }
  const navigate = useNavigate()



  return (
      <>
          <div className="  bg-slate-200 dark:bg-gray-950      h-full w-full max-w-sm  justify-start items-start  max-h-[90vh]  flex flex-col gap-2 ">
              <div
                  id="k-agent-bar"
                  className="flex flex-col gap-2 max-h-[90vh] overflow-y-scroll mt-2 "
              >
                  {agents?.map((agent) => {
                      return (
                          <div
                              key={agent.id}
                              onClick={() => {
                                setOpen(true)
                              }}
                              className={`rounded-sm flex flex-col justify-start items-start gap-2 hover:dark:bg-gray-700 hover:bg-gray-400 cursor-pointer p-2 ${
                                  selected_agent?.id == agent.id &&
                                  ' bg-slate-400 dark:bg-slate-800'
                              }`}
                          >
                              <span className="text-base text-slate-950 dark:text-slate-200">
                                  {agent.name}
                              </span>
                              <span className="text-sm text-slate-500 dark:text-slate-400">
                                  {agent.description}
                              </span>
                          </div>
                      )
                  })}
              </div>
          </div>
          <Modal
              size="medium"
              visible={open}
              onDismiss={() => setOpen(false)}
              header="Not available"
          >
              <Flex className="flex-col gap-2">
                  <span>
                      {' '}
                      This feature is only available on commercial version.
                  </span>
                  <Button
                      onClick={() => {
                          setOpenCal(true)
                      }}
                  >
                      Contact us
                  </Button>
              </Flex>
          </Modal>
          <Modal
              size="large"
              visible={openCal}
              onDismiss={() => setOpenCal(false)}
              header="Not available"
          >
              <Cal
                  namespace="try-enterprise"
                  calLink="team/clearcompass/try-enterprise"
                  style={{
                      width: '100%',
                      height: '100%',
                      overflow: 'scroll',
                  }}
                  config={{ layout: 'month_view' }}
              />
          </Modal>
      </>
  )
}

export default Agents;
