import React,{ useState } from 'react';
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip'

function ArticleCardActions({ citation, addCitation, removeCitation }) {
    const [ anchorEl, setAnchorEl ] = useState(null);
    
    const handleClick = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    return <>
        <Tooltip title="Add Citation">
            <span>
                <IconButton aria-label="Add Citation" onClick={handleClick} disabled={addCitation > -1}>
                </IconButton>
            </span>
        </Tooltip>
        <Tooltip title="Remove Citation">
            <span>
                <IconButton disabled={addCitation > -1} aria-label="Remove Citation" onClick={(e) => removeCitation(e, citation.id)}>
                    <RemoveIcon />
                </IconButton>
            </span>
        </Tooltip>
    </>
}

export default ArticleCardActions;