import { Link, Stack, Typography } from '@mui/material';
import { URLRelation } from './types';

const HistoryEntry = ({ relation }: { relation: URLRelation }) => {
  return (
    <Stack direction={{ xs: 'column', sm: 'column', md: 'row' }} margin={2} spacing={2}
           alignItems={{ xs: 'flex-start', sm: 'flex-start', md: 'center' }}
           justifyContent='space-between'>
      <Typography>{relation.longUrl}</Typography>
      <Link sx={{ fontFamily: 'sans-serif' }} rel='noopener noreferrer' target='_blank'
            href={relation.shortUrl}>{relation.shortUrl}</Link>
    </Stack>
  );
};

export default HistoryEntry;