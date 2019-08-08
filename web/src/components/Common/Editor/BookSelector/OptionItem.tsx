import React from 'react';
import classnames from 'classnames';

import { Option } from '../../../../libs/select';
import CheckIcon from '../../../Icons/Check';
import styles from './OptionItem.scss';

interface Props {
  option: Option;
  isSelected: boolean;
  isFocused: boolean;
  isNew?: boolean;
  onSelect: (Option) => void;
}

function renderBody(isNew: boolean, label: string) {
  if (isNew) {
    return `Create book '${label}'`;
  }

  return label;
}

const OptionItem: React.SFC<Props> = ({
  option,
  isSelected,
  isFocused,
  isNew,
  onSelect
}) => {
  return (
    <button
      role="option"
      type="button"
      aria-selected={isSelected}
      onClick={() => {
        onSelect(option);
      }}
      className={classnames(
        'button-no-ui',
        `book-item-${option.value}`,
        styles['combobox-option'],
        {
          [styles.active]: isSelected,
          [styles.focused]: isFocused
        }
      )}
    >
      {isSelected && (
        <CheckIcon
          fill="white"
          width={12}
          height={12}
          className={styles['check-icon']}
        />
      )}
      <div className={styles['option-label']}>
        {renderBody(isNew, option.label)}
      </div>
    </button>
  );
};

export default OptionItem;
