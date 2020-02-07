/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import classnames from 'classnames';
import { NoteData } from 'jslib/operations/types';
import React, { useState } from 'react';
import { copyToClipboard, selectTextInputValue } from 'web/libs/dom';
import operations from 'web/libs/operations';
import { getNotePath } from 'web/libs/paths';
import { useDispatch } from '../../../store';
import { receiveNote } from '../../../store/note';
import Button from '../../Common/Button';
import Flash from '../../Common/Flash';
import Modal, { Body, Header } from '../../Common/Modal';
import Toggle, { ToggleKind } from '../../Common/Toggle';
import CopyButton from './CopyButton';
import styles from './ShareModal.scss';

// getNoteURL returns the absolute URL for the note
function getNoteURL(uuid: string): string {
  const loc = getNotePath(uuid);
  const path = loc.pathname;

  // TODO: maybe get these values from the configuration instead of parsing
  // current URL.
  const { protocol, host } = window.location;

  return `${protocol}//${host}${path}`;
}

function getHelpText(isPublic: boolean): string {
  if (isPublic) {
    return 'Anyone with this URL can view this note.';
  }

  return 'This note is private only to you.';
}

interface Props {
  isOpen: boolean;
  onDismiss: () => void;
  note: NoteData;
}

const ShareModal: React.FunctionComponent<Props> = ({
  isOpen,
  onDismiss,
  note
}) => {
  const [inProgress, setInProgress] = useState(false);
  const [errMessage, setErrMessage] = useState('');
  const [copyHot, setCopyHot] = useState(false);
  const dispatch = useDispatch();

  const labelId = 'share-note-modal-label';
  const linkValue = getNoteURL(note.uuid);

  function handleToggle(val: boolean) {
    setInProgress(true);

    operations.notes
      .update(note.uuid, { public: val })
      .then(resp => {
        dispatch(receiveNote(resp));
        setInProgress(false);
      })
      .catch(err => {
        console.log('Error sharing note', err);
        setInProgress(false);
        setErrMessage(err.message);
      });
  }

  function handleCopy() {
    copyToClipboard(linkValue);
    setCopyHot(true);

    window.setTimeout(() => {
      setCopyHot(false);
    }, 1200);
  }

  return (
    <Modal
      modalId="T-share-note-modal"
      isOpen={isOpen}
      onDismiss={onDismiss}
      ariaLabelledBy={labelId}
      size="regular"
    >
      <Header
        labelId={labelId}
        heading="Share this note"
        onDismiss={onDismiss}
      />

      <Flash
        kind="danger"
        onDismiss={() => {
          setErrMessage('');
        }}
        hasBorder={false}
        when={Boolean(errMessage)}
        noMargin
      >
        {errMessage}
      </Flash>

      <Body>
        <div className={styles['label-row']}>
          <label htmlFor="link-value" className={styles.label}>
            Link sharing
          </label>

          <Toggle
            id="T-note-public-toggle"
            kind={ToggleKind.green}
            checked={note.public}
            onChange={handleToggle}
            disabled={inProgress}
            label={
              <span className={styles.status}>
                {note.public ? 'Enabled' : 'Disabled'}
              </span>
            }
          />
        </div>

        <input
          id="link-value"
          type="text"
          disabled={!note.public}
          value={linkValue}
          onChange={e => {
            e.preventDefault();
          }}
          className={classnames(
            'form-control text-input text-input-small',
            styles['link-input']
          )}
          onFocus={e => {
            const el = e.target;

            selectTextInputValue(el);
          }}
          onKeyDown={e => {
            e.preventDefault();
          }}
          onMouseUp={e => {
            e.preventDefault();

            const el = e.target as HTMLInputElement;
            selectTextInputValue(el);
          }}
        />

        <p className={styles.help}>{getHelpText(note.public)} </p>

        <div className={styles.actions}>
          {note.public && (
            <CopyButton
              kind="third"
              size="normal"
              onClick={handleCopy}
              isHot={copyHot}
              className={classnames('button-normal button-second', styles.copy)}
            />
          )}

          <Button
            id="T-share-note-modal-close"
            type="button"
            size="normal"
            kind="second"
            onClick={onDismiss}
            disabled={inProgress}
          >
            Cancel
          </Button>
        </div>
      </Body>
    </Modal>
  );
};

export default ShareModal;
