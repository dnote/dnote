/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import React from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import services from 'web/libs/services';
import { BookDomain, RepetitionRuleData } from 'jslib/operations/types';
import { booksToOptions } from 'jslib/helpers/select';
import { getRepetitionsPath, repetitionsPathDef } from 'web/libs/paths';
import Form, { FormState, serializeFormState } from '../Form';
import { useDispatch } from '../../../store';
import { setMessage } from '../../../store/ui';

interface Props extends RouteComponentProps {
  setErrMsg: (string) => void;
  data: RepetitionRuleData;
}

const RepetitionEditContent: React.SFC<Props> = ({
  history,
  setErrMsg,
  data
}) => {
  const dispatch = useDispatch();

  async function handleSubmit(state: FormState) {
    const payload = serializeFormState(state);

    try {
      await services.repetitionRules.update(data.uuid, payload);

      const dest = getRepetitionsPath();
      history.push(dest);

      dispatch(
        setMessage({
          message: `Updated the repetition rule: "${data.title}"`,
          kind: 'info',
          path: repetitionsPathDef
        })
      );
    } catch (e) {
      console.log(e);
      setErrMsg(e.message);
    }
  }

  const initialFormState = {
    title: data.title,
    enabled: data.enabled,
    hour: data.hour,
    minute: data.minute,
    frequency: data.frequency,
    noteCount: data.noteCount,
    bookDomain: data.bookDomain,
    books: booksToOptions(data.books)
  };

  return (
    <Form
      isEditing
      onSubmit={handleSubmit}
      setErrMsg={setErrMsg}
      initialState={initialFormState}
    />
  );
};

export default withRouter(RepetitionEditContent);
