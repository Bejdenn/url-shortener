import { Link, Stack, Typography } from '@mui/material';
import { URLRelation } from './types';

function truncate(value: string, threshold: number) {
  if (value.length > threshold) {
    return value.substring(0, threshold) + '...';
  }

  return value;
}

const HistoryEntry = ({ relation }: { relation: URLRelation }) => {
  return (
    <Stack direction={{ xs: 'column', sm: 'column', md: 'row' }} margin={2} spacing={2}
           alignItems={{ xs: 'flex-start', sm: 'flex-start', md: 'center' }}
           justifyContent='space-between'>
      <Typography>{truncate(relation.longUrl, 65)}</Typography>
      <Link sx={{ fontFamily: 'sans-serif' }} rel='noopener noreferrer' target='_blank'
            href={relation.shortUrl}>{relation.shortUrl}</Link>
    </Stack>
  );
};

export default HistoryEntry;