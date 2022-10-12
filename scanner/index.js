// Copyright (c) 2022 EPAM Systems, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

import fetch from 'node-fetch';
import qs from 'qs';
import dayjs from 'dayjs';

const {
  TARGET_STATUSES: targetStatuses = 'deployed;incomplete',
  STATE_FUNCTION_URL: stateFunctionUrl = 'https://us-central1-superhub.cloudfunctions.net/stacks',
  DAYS_BEFORE: daysBefore = '7',
  DISPOSER_URL: disposerUrl = 'https://stack-disposer-mvn4dxj74a-uc.a.run.app',
  VERBOSE: verbose = 'false',
} = process.env;

/**
 * Get stacks via hub state cloud function
 *
 * @param {object} filter Object contains fields for query
 * @return {array} Array of stack IDs
 */
async function getStacks(filter) {
  const url = `${stateFunctionUrl}?${qs.stringify(filter)}`;
  const body = await fetch(url);
  const data = await body.json();
  return data
      .filter(({id, sandbox}) => id !== 'unset' && sandbox.dir !== 'unset' )
      .map(({id, sandbox}) => ({id, sandbox}));
}

const FORMAT = 'YYYY-MM-DD';

export const scan = async (_, res) => {
  const date = dayjs().subtract(Number.parseInt(daysBefore, 10), 'day');
  // eslint-disable-next-line max-len
  console.log(`Start scanning for stacks in status ${targetStatuses} before ${date.format(FORMAT)}`);
  const stacks = await Promise.all(
      targetStatuses.split(';').map((status) => getStacks({
        status,
        'latestOperation.timestamp[before]': date.format(FORMAT),
      })),
  );

  stacks
      .reduce((acc, value) => ([...acc, ...value]), [])
      .forEach(({id, sandbox: {dir, commit}}) => {
        const params = qs.stringify({commit, verbose});
        const url = `${disposerUrl}/${dir}/${id}?${params}`;
        console.log(`Request undeploy of "${id}" stack`);
        fetch(url, {method: 'DELETE'})
            .then((resp) => {
              const status = resp.ok ? 'undeployed' : 'failed to undeploy';
              console.log(`Stack "${id}" is ${status}`);
            })
            .catch((error) => {
              console.log(`Failed to undeploy stack "${id}"`);
              console.log(error);
            });
      });

  console.log('Finished scanning');

  res.sendStatus(202);
};

