import { Box } from '@mui/material';
import { ReactChild } from 'react';

const Layout = ({ children }: { children: ReactChild[] }) => {
  return (
    <Box sx={{ paddingTop: 4, paddingLeft: 2, paddingRight: 2 }}>
      <Box display={'flex'} justifyContent='center' alignItems='center'>
        <Box
          sx={{
            paddingTop: '6rem',
            maxWidth: 'xl'
          }}
        >{children}</Box>
      </Box>
    </Box>);
};

export default Layout;