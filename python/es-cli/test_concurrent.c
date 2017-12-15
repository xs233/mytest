#include <unistd.h>
#include <stdio.h>
#include <Python.h>

int main(int argc,char** argv)
{
    Py_Initialize();
    if (!Py_IsInitialized())
	{
		perror("Init failed!");
        return -1;
    }
    
	PyRun_SimpleString("import sys");
	PyRun_SimpleString("sys.path.append('./')");

    PyObject* moduleName = PyString_FromString("scroll");
    PyObject* module = PyImport_Import(moduleName);
    if (!module)
	{
		perror("import error");
        return -2;
    }

    PyObject* fun = PyObject_GetAttrString(module, "get_scroll_id");
    if (!fun || !PyCallable_Check(fun))
	{
		perror("get fun failed");
        return -3;
    }

    PyObject* ret = PyObject_CallObject(fun,NULL);
    char* res_scroll = NULL;
    if (ret && PyArg_ParseTuple(ret,"z",res_scroll))
        printf("%s\n",res_scroll);

    if (moduleName)
        Py_DECREF(moduleName);
    if (module)
        Py_DECREF(module);
    if (fun)
        Py_DECREF(fun);
    if (ret)
        Py_DECREF(ret);
    Py_Finalize();

/*
    int i = 0;
    for ( ;i<2;++i)
    {
        pid_t pid = fork();
        if (pid < 0)
            perror("fork failed!\n");
        else if (0 == pid)
            printf("child process : %d\n",getpid());
        else
            printf("father process : %d\n",getpid());
    }
*/
    return 0;
}

