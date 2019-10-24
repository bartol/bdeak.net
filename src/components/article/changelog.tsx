import React from 'react'
import { ManualAutolink } from './manualAutolink'

export const Changelog = ({ changelog }) => {
  return (
    <div className='markdown'>
      <h3 id='changelog'>
        <ManualAutolink id='changelog' />
        Changelog
      </h3>
      <ul>
        {changelog.map(change => {
          const { date, message } = change
          return (
            <li>
              {date} - {message}
            </li>
          )
        })}
      </ul>
    </div>
  )
}