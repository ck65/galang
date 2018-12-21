function switchTab(tabid)
{
    if(tabid == '') return false;
    s = document.getElementsByTagName("tr");
    console.log(s)
    for(let i=1; i<=s.length-1; i++)
    {
        if(tabid == 't_'+i){
            q = document.getElementsByName('tl_'+i);
            console.log(q)
            for(let w=0;w<q.length;w++){
                q[w].hidden = false
            }
        }
        else
        {
            console.log(document.getElementsByName('tl_'+i))
            e = document.getElementsByName('tl_'+i);
            console.log(e)
            for(let w=0;w<e.length;w++){
                e[w].hidden = true
            }
        }

    }
    for(let i=0; i<=9; i++)
    {
        if(tabid == 't_'+i){
            document.getElementById(tabid).style.color="red";
        }
        else
        {
            document.getElementById('t_'+i).style.color="white";
        }

    }
    return true;
}