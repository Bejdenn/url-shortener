import { Link, Stack, Typography } from '@mui/material';
import { URLRelation } from './types';
import { API_URL } from './App';

function truncate(value: string, threshold: number) {
  if (value.length > threshold) {
    return value.substring(0, threshold) + '...';
  }

  return value;
}

interface HistoryEntryProps {
  relation: URLRelation;
  isDev: boolean;
}

const HistoryEntry = ({ relation, isDev }: HistoryEntryProps) => {
  const targetAddr = isDev ? API_URL + '/' + relation.id : relation.shortUrl;

  return (
    <Stack direction={{ xs: 'column', sm: 'column', md: 'row' }} margin={2} spacing={2}
           alignItems={{ xs: 'flex-start', sm: 'flex-start', md: 'center' }}
           justifyContent='space-between'>
      <Typography>{truncate(relation.longUrl, 65)}</Typography>
      <Link sx={{ fontFamily: 'sans-serif' }} rel='noopener noreferrer' target='_blank'
            href={targetAddr}>{targetAddr}</Link>
    </Stack>
  );
};

export default HistoryEntry;