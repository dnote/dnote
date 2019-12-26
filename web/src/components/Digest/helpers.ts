import { DigestData } from 'jslib/operations/types';

// getDigestTitle returns a title for the digest
export function getDigestTitle(digest: DigestData) {
  return `${digest.repetitionRule.title} #${digest.version}`;
}
