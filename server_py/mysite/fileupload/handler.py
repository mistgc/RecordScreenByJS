import os
def handle_uploaded_file(path, f):
    path_dir = path.split('/')
    isExists = os.path.exists('./' + path_dir[1] + '/' + path_dir[2])
    for i in path_dir:
        print(i)
    if not isExists:
        # os.makedirs('./' + path_dir[0])
        os.makedirs( './' + path_dir[1] + '/' + path_dir[2])
    with open('.' + path, 'ab+') as destination:
        for chunk in f.chunks():
            destination.write(chunk)
